package sweetyml

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/go-redis/redis/v8"
	"log"
	"sweet-common/utils"
)

func initRedis() {
	conf := GetYmlConf()
	if !conf.Sweet.RedisConfig.Active {
		return
	}
	logs.Info("Init Redis....")
	addr := fmt.Sprintf("%s:%d", conf.Sweet.RedisConfig.Host, conf.Sweet.RedisConfig.Port)
	pwd := conf.Sweet.RedisConfig.Password

	utils.Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,                            // Redis服务器地址
		Password: pwd,                             // Redis密码，如果没有设置则留空
		DB:       conf.Sweet.RedisConfig.Database, // Redis数据库索引，默认为0
	})
	utils.Ctx = context.Background()
	_, err := utils.Rdb.Ping(utils.Ctx).Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
	}
	fmt.Println("Successfully connected to Redis!")
}
