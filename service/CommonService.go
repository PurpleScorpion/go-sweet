package service

import (
	"context"
	"fmt"
	"github.com/PurpleScorpion/go-sweet-orm/mapper"
	"github.com/beego/beego/orm"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	sweetyml "go-sweet/common/yaml"
	"log"
	"time"
)

var (
	rdb *redis.Client
	ctx context.Context
)

func InitService() {
	initMySQL()
	initRedis()
}

func initMySQL() {
	conf := sweetyml.GetYmlConf()
	if !conf.Sweet.MySqlConfig.Active {
		return
	}
	username := conf.Sweet.MySqlConfig.User
	password := conf.Sweet.MySqlConfig.Password
	host := conf.Sweet.MySqlConfig.Host
	port := conf.Sweet.MySqlConfig.Port
	dbName := conf.Sweet.MySqlConfig.DbName

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", username, password, host, port, dbName)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", connStr)
	orm.SetMaxIdleConns("default", conf.Sweet.MySqlConfig.MaxIdleConns)
	orm.SetMaxOpenConns("default", conf.Sweet.MySqlConfig.MaxOpenConns)
	orm.Debug = false
	mapper.InitMapper(mapper.MySQL, true)
}

func initRedis() {
	conf := sweetyml.GetYmlConf()
	if !conf.Sweet.RedisConfig.Active {
		return
	}

	addr := fmt.Sprintf("%s:%d", conf.Sweet.RedisConfig.Host, conf.Sweet.RedisConfig.Port)
	pwd := conf.Sweet.RedisConfig.Password

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,                            // Redis服务器地址
		Password: pwd,                             // Redis密码，如果没有设置则留空
		DB:       conf.Sweet.RedisConfig.Database, // Redis数据库索引，默认为0
	})
	ctx = context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
	}
	fmt.Println("Successfully connected to Redis!")
}

func SetCache(key string, value string) bool {
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Fatalf("Setting key value pairs failed: %v", err)
		return false
	}
	return true
}

func SetCache4Expiration(key string, value string, expiration time.Duration) bool {
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Fatalf("Setting key value pairs failed: %v", err)
		return false
	}
	return true
}

func DeleteCache(key string) bool {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Fatalf("Deleting key value pairs failed: %v", err)
		return false
	}
	return true
}

func GetCache(key string) string {
	val, err := rdb.Get(ctx, key).Result()
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

func GetMinute(num int) time.Duration {
	return time.Duration(num) * time.Minute
}
