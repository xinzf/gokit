package server

import (
	"fmt"
	"github.com/xinzf/gokit/logger"
	"github.com/xinzf/gokit/registry"
	"github.com/xinzf/gokit/utils"
	"net"
)

type Options struct {
	name     string
	host     string
	port     int
	registry registry.Registry
	logger   logger.Logger
}

func newOptions(opt ...Option) Options {
	var opts Options

	if len(opt) > 0 {
		for _, o := range opt {
			o(&opts)
		}
	}

	//if opts.registry == nil {
	//	opts.registry = registry.DefaultRegistry
	//}
	if opts.host == "" {
		opts.host = localIP
	}
	if opts.port == 0 {
		l, _ := net.Listen("tcp", ":0")
		opts.port = l.Addr().(*net.TCPAddr).Port
		l.Close()
	}
	if opts.name == "" {
		opts.name = fmt.Sprintf("gokit-%s-%s", opts.host, utils.UUID())
	}
	if opts.logger == nil {
		opts.logger = logger.DefaultLogger
	}

	return opts
}

type Option func(options *Options)

func Name(name string) Option {
	return func(options *Options) {
		options.name = name
	}
}

func Host(host string) Option {
	return func(options *Options) {
		options.host = host
	}
}

func Port(port int) Option {
	return func(options *Options) {
		options.port = port
	}
}

func Registry(r registry.Registry) Option {
	return func(options *Options) {
		options.registry = r
	}
}

func Logger(logger logger.Logger) Option {
	return func(options *Options) {
		options.logger = logger
	}
}
