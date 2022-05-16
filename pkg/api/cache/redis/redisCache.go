package redis

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"time"
)

type CacheRedis struct {
	Host    string
	ExpTime time.Duration
	Ctx     context.Context
}

func InitializeRedisCache(host string, expTime time.Duration, ctx context.Context) *CacheRedis {
	return &CacheRedis{
		Host:    host,
		ExpTime: expTime,
		Ctx:     ctx,
	}
}

func (cache *CacheRedis) redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cache.Host,
	})
}

func (cache *CacheRedis) Set(key string, value interface{}) error {

	client := cache.redisClient()

	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Err(err)
		return err
	}

	return client.Set(cache.Ctx, key, jsonValue, cache.ExpTime).Err()

}
func (cache *CacheRedis) Get(key string) (value interface{}) {

	client := cache.redisClient()

	getValue, err := client.Get(cache.Ctx, key).Bytes()
	if err != nil {
		log.Err(err)
		return err
	}

	var returnValue interface{}

	err = json.Unmarshal(getValue, &returnValue)
	if err != nil {
		log.Err(err)
		return err
	}

	return returnValue
}

func (cache *CacheRedis) Delete(key string) error {
	client := cache.redisClient()
	return client.Del(cache.Ctx, key).Err()
}
