package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xinzf/gokit/registry"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	_runtime "runtime"
	"strings"
	"syscall"
	"time"
)

const (
	ReqMetaDataKey = "__share_data"
)

var srv *server

func New(opt ...Option) {
	_runtime.GOMAXPROCS(_runtime.NumCPU())

	opts := newOptions(opt...)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders(opts.allowHeaders...)

	e := gin.New()
	e.Use(gin.Recovery(), gin.Logger(), cors.New(config))
	e.MaxMultipartMemory = 8 << 20 // 8 MiB

	srv = &server{
		id:       fmt.Sprintf("%s-%s-%d", hostName, opts.name, time.Now().UnixNano()),
		options:  newOptions(opt...),
		handlers: map[string]map[string]*Method{},
		metaData: map[string]string{},
		g:        e,
	}
}

type server struct {
	id       string
	options  Options
	handlers map[string]map[string]*Method
	metaData map[string]string
	g        *gin.Engine
}

type Initialization func() error

func Register(name string, hdl interface{}) {
	refVal := reflect.ValueOf(hdl)
	refType := reflect.TypeOf(hdl)
	hdlName := refType.Elem().String()

	if srv.handlers[name] == nil {
		srv.handlers[name] = make(map[string]*Method)
	}

	for i := 0; i < refType.NumMethod(); i++ {
		m, err := newMethod(hdlName, refType.Method(i).Name, refVal.Method(i))
		if err != nil {
			srv.options.logger.Warn(err)
			continue
		}

		srv.handlers[name][refType.Method(i).Name] = m
	}
}

func Run(ctx context.Context, initFn ...Initialization) {
	for _, fn := range initFn {
		if err := fn(); err != nil {
			srv.options.logger.Fatal(err)
			return
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	if err := srv.register(); err != nil {
		srv.options.logger.Fatal(err)
	}

	go srv.serve(fmt.Sprintf("%s:%d", srv.options.host, srv.options.port))

	<-ch

	Stop()
}

func (this *server) serve(addr string) {
	this.options.logger.Info(fmt.Sprintf("server listen: %s", fmt.Sprintf("%s:%d", this.options.host, this.options.port)))

	this.g.Use(gin.Recovery())

	this.g.NoRoute(this.call)

	err := http.ListenAndServe(addr, this.g)
	if err != nil {
		this.options.logger.Fatal(err)
		return
	}
}

func (this *server) register() error {
	if this.options.registry == nil {
		return nil
	}

	for hdlName, mp := range this.handlers {
		for methodName, method := range mp {
			hdlName = strings.Replace(hdlName, ".", "_", -1)
			this.metaData[fmt.Sprintf("%s_%s", hdlName, methodName)] = method.Structure()
		}
	}

	err := this.options.registry.Register(&registry.Service{
		ID:      this.id,
		Name:    this.options.name,
		Address: this.options.host,
		Port:    this.options.port,
		Metas:   this.metaData,
	})

	if err != nil {
		this.options.logger.Fatal(fmt.Sprintf("service: %s register failed,err: %s", this.options.name, err.Error()))
		return err
	}

	this.options.logger.Info("register server success")
	return nil
}

func Stop() {
	if srv.options.registry != nil {
		if err := srv.options.registry.Deregister(srv.id); err != nil {
			srv.options.logger.Error(fmt.Sprintf("deregister server failed, err: %s", err.Error()))
		} else {
			srv.options.logger.Info("deregister server success")
		}
	}

	srv.options.logger.Info("server stopped.")
}

func (this *server) encodeMetadata(md map[string]string) []string {
	var tags []string
	for k, v := range md {
		tags = append(tags, fmt.Sprintf("%s:%s", k, v))
	}

	return tags
}

func (this *server) errHandler(err error, code ...int) interface{} {
	var c = -1000
	if len(code) > 0 {
		c = code[0]
	}
	mp := map[string]interface{}{
		"code":    c,
		"message": err.Error(),
		"result":  map[string]interface{}{},
	}
	return mp
}

func (this *server) call(ctx *gin.Context) {

	if ctx.Request.Method == http.MethodHead {
		ctx.JSON(200, nil)
		return
	}

	name := ctx.Query("_service")
	method := ctx.Query("_method")
	if name == "" {
		ctx.JSON(200, this.errHandler(errors.New("missing service name")))
		return
	}
	if method == "" {
		ctx.JSON(200, this.errHandler(errors.New("missing method name")))
		return
	}

	methods, found := this.handlers[name]
	if !found {
		ctx.JSON(200, this.errHandler(fmt.Errorf("gokit server: %s not found", name)))
		return
	}

	fn, found := methods[method]
	if !found {
		ctx.JSON(200, this.errHandler(fmt.Errorf("method: %s.%s not found", name, method)))
		return
	}

	var (
		rawData []byte
		err     error
	)
	rawData = make([]byte, 0)

	var rsp interface{}
	c := context.WithValue(context.Background(), ReqMetaDataKey, ctx.Request.Header)
	if ctx.ContentType() == gin.MIMEJSON {
		rawData, err = ctx.GetRawData()
		if err != nil {
			ctx.JSON(200, this.errHandler(fmt.Errorf("failed to get request data,err: %s", err.Error())))
			return
		}
		rsp, err = fn.Call(c, rawData)
	} else if ctx.ContentType() == gin.MIMEMultipartPOSTForm {
		rsp, err = fn.CallUpload(c, ctx)
	}
	if err != nil {
		ctx.JSON(200, this.errHandler(err))
		return
	}
	ctx.JSON(200, rsp)
}
