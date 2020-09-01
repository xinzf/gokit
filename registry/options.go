package registry

import (
	"context"
	"time"
)

type Options struct {
	Addrs    []string
	Timeout  time.Duration
	Context  context.Context
	Interval time.Duration
}

type Option func(*Options)

func newOptions(opt ...Option) Options {
	opts := Options{
		Addrs:    []string{"http://127.0.0.1:8500"},
		Timeout:  3 * time.Second,
		Context:  context.Background(),
		Interval: 5 * time.Second,
	}

	if len(opt) > 0 {
		for _, o := range opt {
			o(&opts)
		}
	}

	return opts
}

// Addrs is the registry addresses to use
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

func Interval(t time.Duration) Option {
	return func(o *Options) {
		o.Interval = t
	}
}

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}
