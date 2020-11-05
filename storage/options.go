package storage

import (
	"fmt"
	"github.com/xinzf/gokit/logger"
)

type DbOptions struct {
	Addr         string
	User         string
	Pswd         string
	Name         string
	Log          bool
	MaxIdleConns int
	MaxOpenConns int
	initSql      []string
}

type DbOption func(o *DbOptions)

func newDbOptions(opt ...DbOption) DbOptions {
	opts := DbOptions{
		Addr:         "127.0.0.1",
		User:         "root",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		initSql:      []string{},
	}

	if len(opt) > 0 {
		for _, o := range opt {
			o(&opts)
		}
	}
	return opts
}

func DbHost(addr string) DbOption {
	return func(o *DbOptions) {
		o.Addr = addr
	}
}

func DbUser(user string) DbOption {
	return func(o *DbOptions) {
		o.User = user
	}
}

func DbPassword(pswd string) DbOption {
	return func(o *DbOptions) {
		o.Pswd = pswd
	}
}

func DbName(name string) DbOption {
	return func(o *DbOptions) {
		o.Name = name
	}
}

func DbLog() DbOption {
	return func(o *DbOptions) {
		o.Log = true
	}
}

func DbMaxIdleConns(num int) DbOption {
	return func(o *DbOptions) {
		o.MaxIdleConns = num
	}
}

func DbMaxOpenConns(num int) DbOption {
	return func(o *DbOptions) {
		o.MaxOpenConns = num
	}
}

func InitQuery(sql string) DbOption {
	return func(o *DbOptions) {
		o.initSql = append(o.initSql, sql)
	}
}

func (s DbOptions) String() string {
	u := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true&parseTime=true&loc=Local",
		s.User,
		s.Pswd,
		s.Addr,
		s.Name)
	return u
}

type MongoOptions struct {
	Addr   string `mapstructure:"addr"`
	Debug  bool   `mapstructure:"debug"`
	Name   string `mapstructure:"name"`
	logger logger.Logger
}

type MongoOption func(options *MongoOptions)

func newMongoOptions(opt ...MongoOption) MongoOptions {
	opts := MongoOptions{
		Addr:   "127.0.0.1:27017",
		logger: logger.DefaultLogger,
	}

	if len(opt) > 0 {
		for _, o := range opt {
			o(&opts)
		}
	}
	return opts
}

func MongoAddr(addr, name string) MongoOption {
	return func(o *MongoOptions) {
		o.Addr = addr
		o.Name = name
	}
}

func MongoDebug(flag bool) MongoOption {
	return func(o *MongoOptions) {
		o.Debug = flag
	}
}

func MongoLogger(logger logger.Logger) MongoOption {
	return func(options *MongoOptions) {
		options.logger = logger
	}
}

type RedisOptions struct {
	Addr   string `mapstructure:"addr"`
	Pswd   string `mapstructure:"pswd"`
	DB     int    `mapstructure:"db"`
	Logger logger.Logger
}

func newRedisOptions(opt ...RedisOption) RedisOptions {
	opts := RedisOptions{
		Addr:   "127.0.0.1:6379",
		Logger: logger.DefaultLogger,
	}

	if len(opt) > 0 {
		for _, o := range opt {
			o(&opts)
		}
	}

	return opts
}

type RedisOption func(o *RedisOptions)

func RedisConfig(addr, pswd string, db int) RedisOption {
	return func(o *RedisOptions) {
		o.Addr = addr
		o.Pswd = pswd
		o.DB = db
	}
}

func RedisLogger(logger logger.Logger) RedisOption {
	return func(o *RedisOptions) {
		o.Logger = logger
	}
}
