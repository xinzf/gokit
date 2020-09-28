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
	path            string
	lastModifyIndex map[string]uint64
	stopChan        chan string
	ctx             context.Context
	cancle          context.CancelFunc

	handler func(result *Result)
}

func newWatcher(path string, hdl func(result *Result)) *consulWatcher {
	return &consulWatcher{
		path:            path,
		lastModifyIndex: map[string]uint64{},
		handler:         hdl,
	}
}

func (this *consulWatcher) isModified(path string, index uint64) bool {
	lastIndex, found := this.lastModifyIndex[path]
	if !found {
		this.lastModifyIndex[path] = index
		return true
	}
	if lastIndex < index {
		this.lastModifyIndex[path] = index
		return true
	}

	return false
}

func (this *consulWatcher) run(ctx context.Context, address string, stopChan chan string) error {
	var err error
	this.Plan, err = watch.Parse(map[string]interface{}{
		"type":   "keyprefix",
		"prefix": this.path,
	})
	if err != nil {
		return err
	}

	this.stopChan = stopChan
	this.ctx, this.cancle = context.WithCancel(ctx)

	this.Plan.HybridHandler = func(val watch.BlockingParamVal, data interface{}) {
		kvPairs := data.(api.KVPairs)
		for _, kv := range kvPairs {
			path := strings.TrimSuffix(strings.TrimPrefix(kv.Key, "/"), "/")
			if this.isModified(path, kv.ModifyIndex) {
				this.handler(&Result{
					Path:        path,
					Value:       kv.Value,
					ModifyIndex: kv.ModifyIndex,
				})
			}
		}
	}

	errChan := make(chan error)
	go func() {
		errChan <- this.RunWithConfig(address, &api.Config{
			Address:   address,
			Transport: this.newTransport(),
		})
	}()

	defer func() {
		this.Stop()
		this.stopChan <- this.path

		if err != nil {
			logger.DefaultLogger.Error("watcher", fmt.Sprintf("watcher: %s stop with err: %s", this.path, err.Error()))
		} else {
			logger.DefaultLogger.Info("watcher", fmt.Sprintf("watcher: %s stop", this.path))
		}
	}()

	logger.DefaultLogger.Info("watcher", fmt.Sprintf("watcher: %s start success", this.path))

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
	if path == "" {
		path = this.prefix
	} else {
		path = this.prefix + "/" + path
	}
	watcher, found := this.watchers[path]
	if found {
		watcher.stop()
	}

	this.watchers[path] = newWatcher(path, hdl)
	go this.watchers[path].run(this.ctx, this.address, this.stopChan)
}

func (this *ConsulWatcher) Remove(path string) {
	path = strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
	path = this.prefix + "/" + path
	watcher, found := this.watchers[path]
	if found {
		watcher.stop()
	}
}

func (this *ConsulWatcher) Shutdown() {
	this.cancle()
}

func (this *ConsulWatcher) Run() {
	defer func() {
		logger.DefaultLogger.Info("watcherPool", "shutdown success")
	}()

	logger.DefaultLogger.Info("watcherPool", "run success")
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
