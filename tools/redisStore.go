package tools

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisStore struct {
	redisClient *redis.Client
}

var ctx = context.Background()

// NewRedisClient 构造函数
func (rs *RedisStore) NewRedisClient() *RedisStore {
	// 从全局配置中获取redis配置
	config := GetConfig().RedisConfig

	// 初识化redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})

	redisStore := &RedisStore{redisClient: client}
	return redisStore
}

// Set sets the digits for the captcha id.
func (rs *RedisStore) Set(id string, value string) error {
	err := rs.redisClient.Set(ctx, id, value, time.Minute*1).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (rs *RedisStore) Get(id string, clear bool) string {
	val, err := rs.redisClient.Get(ctx, id).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if clear {
		err := rs.redisClient.Del(ctx, id).Err()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return val
}

// Verify captcha's answer directly
func (rs *RedisStore) Verify(id, answer string, clear bool) bool {
	val := rs.Get(id, clear)
	return val == answer
}
