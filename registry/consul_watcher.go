package registry

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/xinzf/gokit/logger"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type consulWatcher struct {
	*watch.Plan
	prefix          string
	path            string
	lastModifyIndex map[string]uint64
	stopChan        chan string
	ctx             context.Context
	cancle          context.CancelFunc
	consulAddr      string
	//removed         bool
	//lock            *sync.Mutex

	handler func(result *Result)
}

func newWatcher(prefix, path, consulAddr string, hdl func(result *Result)) *consulWatcher {
	return &consulWatcher{
		prefix:          prefix,
		path:            path,
		lastModifyIndex: map[string]uint64{},
		handler:         hdl,
		consulAddr:      consulAddr,
		//lock:            new(sync.Mutex),
	}
}

func (this *consulWatcher) isModified(index uint64) (string, bool) {
	lastIndex, found := this.lastModifyIndex[this.path]
	if !found {
		this.lastModifyIndex[this.path] = index
		return "create", true
	}
	if lastIndex < index {
		this.lastModifyIndex[this.path] = index
		return "modify", true
	}

	return "", false
}

func (this *consulWatcher) getAbsPath() string {
	keys := strings.Split(this.path, "/")
	if len(keys) == 0 {
		return this.prefix
	}

	if len(keys[0]) == 0 {
		return this.prefix
	}

	if len(this.prefix) == 0 {
		return strings.Join(keys, "/")
	}

	return this.prefix + "/" + strings.Join(keys, "/")
}

func (this *consulWatcher) run(ctx context.Context, stopChan chan string) error {
	var err error
	this.Plan, err = watch.Parse(map[string]interface{}{
		"type": "key",
		"key":  this.getAbsPath(),
	})
	if err != nil {
		return err
	}

	this.stopChan = stopChan
	this.ctx, this.cancle = context.WithCancel(ctx)

	this.Plan.HybridHandler = func(val watch.BlockingParamVal, data interface{}) {
		if data == nil {
			this.lastModifyIndex = map[string]uint64{}
			this.handler(&Result{
				Path:        this.path,
				Value:       nil,
				ModifyIndex: 0,
				Event:       "delete",
			})
			return
		}

		kvPair := data.(*api.KVPair)
		if event, is := this.isModified(kvPair.ModifyIndex); is {
			this.handler(&Result{
				Path:        this.path,
				Value:       kvPair.Value,
				ModifyIndex: kvPair.ModifyIndex,
				Event:       event,
			})
		}
	}

	errChan := make(chan error)
	go func() {
		errChan <- this.RunWithConfig(this.consulAddr, &api.Config{
			Address:   this.consulAddr,
			Transport: this.newTransport(),
			WaitTime:  time.Second,
		})
	}()

	defer func() {
		this.Stop()
		this.stopChan <- this.path

		if err != nil {
			logger.DefaultLogger.Error("watcher", fmt.Sprintf("%s stop with err: %s", this.path, err.Error()))
		} else {
			logger.DefaultLogger.Info("watcher", fmt.Sprintf("%s stop", this.path))
		}
	}()

	logger.DefaultLogger.Info("watcher", fmt.Sprintf("%s start success", this.path))

	for {
		select {
		case <-this.ctx.Done():
			return nil
		case err = <-errChan:
			return nil
		}
	}
}

func (this *consulWatcher) newTransport() *http.Transport {
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	runtime.SetFinalizer(&t, func(tr **http.Transport) {
		(*tr).CloseIdleConnections()
	})
	return t
}

func (this *consulWatcher) stop() {
	this.cancle()
}

type ConsulWatcher struct {
	prefix   string
	watchers map[string]*consulWatcher
	ctx      context.Context
	cancle   context.CancelFunc
	address  string
	stopChan chan string
}

func NewConsulWatcher(address, prefix string) *ConsulWatcher {
	prefix = strings.TrimSuffix(strings.TrimPrefix(prefix, "/"), "/")
	ctx, cancel := context.WithCancel(context.Background())
	wtc := &ConsulWatcher{
		prefix:   prefix,
		address:  address,
		watchers: map[string]*consulWatcher{},
		stopChan: make(chan string),
		ctx:      ctx,
		cancle:   cancel,
	}
	DefaultWatcher = wtc
	return wtc
}

func (this *ConsulWatcher) Add(path string, hdl func(result *Result)) {
	path = strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
	watcher, found := this.watchers[path]
	if found {
		watcher.stop()
	}

	this.watchers[path] = newWatcher(this.prefix, path, this.address, hdl)
	go this.watchers[path].run(this.ctx, this.stopChan)
}

func (this *ConsulWatcher) Remove(key ...string) {
	path := strings.TrimPrefix(strings.TrimSuffix(strings.Join(key, "/"), "/"), "/")
	watcher, found := this.watchers[path]
	if found {
		watcher.stop()
		return
	}
	for n, _ := range this.watchers {
		logger.DefaultLogger.Debug(n)
	}
}

func (this *ConsulWatcher) Shutdown() {
	this.cancle()
}

func (this *ConsulWatcher) Run() {
	defer func() {
		logger.DefaultLogger.Info("watcherPool", this.prefix, "shutdown success")
	}()

	logger.DefaultLogger.Info("watcherPool", this.prefix, "run success")
BREAK:
	for {
		select {
		case <-this.ctx.Done():
			break BREAK
		case path := <-this.stopChan:
			delete(this.watchers, path)
		}
	}

	time.Sleep(time.Second)
}
