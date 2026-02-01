package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var (
	Rdb *redis.Client
	Ctx context.Context
)

func SetCache(key string, value string) bool {
	err := Rdb.Set(Ctx, key, value, 0).Err()
	if err != nil {
		log.Fatalf("Setting key value pairs failed: %v", err)
		return false
	}
	return true
}

func SetCache4Expiration(key string, value string, expiration time.Duration) bool {
	err := Rdb.Set(Ctx, key, value, expiration).Err()
	if err != nil {
		log.Fatalf("Setting key value pairs failed: %v", err)
		return false
	}
	return true
}

func DeleteCache(key string) bool {
	err := Rdb.Del(Ctx, key).Err()
	if err != nil {
		log.Fatalf("Deleting key value pairs failed: %v", err)
		return false
	}
	return true
}

// 获取所有以某前缀的缓存
func GetAllCache(prefix string) []string {
	return Rdb.Keys(Ctx, prefix+"*").Val()
}

func GetCache(key string) string {
	val, err := Rdb.Get(Ctx, key).Result()
	if err == redis.Nil {
		return "null"
	} else if err != nil {
		return "null"
	}
	return val
}

func GetHour(num int) time.Duration {
	return time.Duration(num) * time.Hour
}

func GetSecond(num int) time.Duration {
	return time.Duration(num) * time.Second
}

func GetMinute(num int) time.Duration {
	return time.Duration(num) * time.Minute
}

func GetRedisClient() *redis.Client {
	return Rdb
}

func GetRedisCtx() context.Context {
	return Ctx
}
