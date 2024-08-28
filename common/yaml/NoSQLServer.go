package sweetyml

import (
	"context"
	"fmt"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/beego/beego/v2/core/logs"
	"github.com/go-redis/redis/v8"
	"log"
	"sweet-common/utils"
)

func initRedis() {
	host := keqing.ValueString("${sweet.redis.host}")
	if keqing.IsEmpty(host) {
		return
	}
	logs.Info("Init Redis....")

	port := keqing.ValueInt("${sweet.redis.port}")
	if port == 0 {
		port = 6379
	}
	if port <= 0 || port > 65535 {
		panic("Redis port is invalid")
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	password := keqing.ValueString("${sweet.redis.password}")
	database := keqing.ValueInt("${sweet.redis.database}")
	if database < 0 || database > 15 {
		panic("Redis database is invalid")
	}

	utils.Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,     // Redis服务器地址
		Password: password, // Redis密码，如果没有设置则留空
		DB:       database, // Redis数据库索引，默认为0
	})
	utils.Ctx = context.Background()
	_, err := utils.Rdb.Ping(utils.Ctx).Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
	}
}
