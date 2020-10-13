package storage

import (
	"errors"
	"github.com/go-redis/redis"
	"sync"
)

var Redis *_redis
var redisOnce sync.Once

type _redis struct {
	options RedisOptions
	client  *redis.Client
	once    sync.Once
}

func (this *_redis) Init(opt ...RedisOption) error {

	var err error
	redisOnce.Do(func() {
		Redis = &_redis{options: newRedisOptions(opt...)}
		if Redis.options.Addr == "" {
			err = errors.New("redis config's addr is empty")
			return
		}

		Redis.client = redis.NewClient(&redis.Options{
			Addr:     Redis.options.Addr,
			Password: Redis.options.Pswd,
			DB:       Redis.options.DB,
		})

		_, err = Redis.client.Ping().Result()
		if err != nil {
			err = errors.New("redis connect failed,err: " + err.Error())
			return
		}
		Redis.options.Logger.Info("Redis", "Redis inited.")
	})

	return err
}

func (this *_redis) Client() *redis.Client {
	return this.client
}
