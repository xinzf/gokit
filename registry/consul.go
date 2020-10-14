package registry

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type _consulRegistry struct {
	options   Options
	client    *api.Client
	services  map[string][]*Service
	connected bool

	sync.RWMutex
	once    sync.Once
	isWatch bool
}

func NewConsul(opt ...Option) Registry {
	consul := &_consulRegistry{
		options:  newOptions(opt...),
		services: map[string][]*Service{},
	}

	consul.Init()

	if err := consul.connect(); err != nil {
		panic(err)
	}

	DefaultRegistry = consul
	return consul
}

func (this *_consulRegistry) Init(opt ...Option) {
	if len(opt) > 0 {
		this.options = newOptions(opt...)
	}

	if err := this.connect(); err != nil {
		panic(err)
	}
}

func (this *_consulRegistry) Register(svc *Service) error {
	asr := &api.AgentServiceRegistration{
		ID:      svc.ID,
		Name:    svc.Name,
		Port:    svc.Port,
		Address: svc.Address,
		Meta:    svc.Metas,
		Check: &api.AgentServiceCheck{
			Interval:                       this.options.Interval.String(),
			Timeout:                        this.options.Timeout.String(),
			HTTP:                           svc.Domain(),
			Method:                         "HEAD",
			Notes:                          fmt.Sprintf("Check server: %s", svc.Name),
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	if err := this.client.Agent().ServiceRegister(asr); err != nil {
		return err
	}

	return nil
}

func (this *_consulRegistry) Deregister(id string) error {
	return this.client.Agent().ServiceDeregister(id)
}

func (this *_consulRegistry) Watch(ctx context.Context) error {
	if err := this.loadServices(); err != nil {
		this.connected = false
		return err
	}
	this.isWatch = true

	this.once.Do(func() {
		go func() {
			ticker := time.NewTicker(3 * time.Second)
			defer func() {
				ticker.Stop()
			}()

			for {
				select {
				case <-ticker.C:
					if err := this.loadServices(); err != nil {
						this.connected = false
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	})

	return nil
}

func (this *_consulRegistry) connect() error {
	config := api.DefaultConfig()

	if len(this.options.Addrs) == 0 {
		return errors.New("missing consul address")
	}
	config.Address = this.options.Addrs[0]

	if config.HttpClient == nil {
		config.HttpClient = new(http.Client)
	}
	config.HttpClient.Transport = this.newTransport()

	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	_, err = client.Agent().Host()
	if err != nil {
		return err
	}

	this.client = client
	this.connected = true
	return nil
}

func (this *_consulRegistry) newTransport() *http.Transport {
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

func (this *_consulRegistry) loadServices() error {
	if this.connected == false {
		if err := this.connect(); err != nil {
			return err
		}
	}

	rsp, _, err := this.client.Catalog().Services(&api.QueryOptions{})
	if err != nil {
		return err
	}

	services := make(map[string][]*Service)
	for name, _ := range rsp {
		resp, _, err := this.client.Health().Service(name, "", true, &api.QueryOptions{})
		if err != nil {
			return err
		}

		svcs := make([]*Service, 0)
		for _, s := range resp {
			if s.Checks.AggregatedStatus() != "passing" {
				continue
			}

			svcs = append(svcs, &Service{
				ID:      s.Service.ID,
				Name:    s.Service.Service,
				Address: s.Service.Address,
				Port:    s.Service.Port,
				Metas:   s.Service.Meta,
			})
		}

		if len(svcs) == 0 {
			continue
		}

		services[name] = svcs
	}

	this.Lock()
	this.services = services
	this.Unlock()

	return nil
}

func (this *_consulRegistry) Services() map[string][]*Service {
	if !this.isWatch {
		this.loadServices()
	}
	defer this.RUnlock()
	this.RLock()
	return this.services
}

func (this *_consulRegistry) Service(name string) (*Service, bool) {
	this.RLock()
	defer this.RUnlock()

	svcs, found := this.services[name]
	if !found {
		return nil, false
	}

	if len(svcs) == 0 {
		return nil, false
	}

	rand.Seed(time.Now().UnixNano())
	svc := svcs[rand.Intn(len(svcs))]

	return svc, true
}
