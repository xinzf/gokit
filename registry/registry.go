package registry

import "context"

type Registry interface {
	Init(opt ...Option)
	Register(svc *Service) error
	Deregister(id string) error
	Watch(ctx context.Context) error
	Services() map[string][]*Service
	Service(name string) (*Service, bool)
}

var DefaultRegistry Registry
