package registry

import (
	"context"
	jsoniter "github.com/json-iterator/go"
)

type Registry interface {
	Init(opt ...Option)
	Register(svc *Service) error
	Deregister(id string) error
	Watch(ctx context.Context) error
	Services() map[string][]*Service
	Service(name string) (*Service, bool)
}

type Watcher interface {
	Add(path string, hdl func(result *Result))
	Remove(path string)
	Shutdown()
	Run()
}

type Result struct {
	Path        string
	Value       []byte
	ModifyIndex uint64
}

func (this *Result) String() string {
	return string(this.Value)
}

func (this *Result) Bind(obj interface{}) error {
	return jsoniter.Unmarshal(this.Value, obj)
}

var DefaultRegistry Registry
var DefaultWatcher Watcher
