package registry

type Registry interface {
	Init(opt ...Option)
	Register(svc *Service) error
	Deregister(id string) error
	Watch()
	Services() map[string][]*Service
	Service(name string) (*Service, bool)
}

//var (
//	DefaultRegistry = newConsul()
//)
