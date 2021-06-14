package services

import (
	"Panda/common/db"
	"Panda/common/log"
	queue "Panda/common/queue/redis"
	"Panda/common/response"
	"Panda/conf"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"time"
)

var (
	DefaultDBClient    *gorm.DB
	DefaultRedisClient *redis.Client
	RedisQueue         *queue.Queue
	respError          response.RespError
	NewCode            = func(code int) response.RespError { return response.NewCode(code) }
	NewError           = func(err interface{}) response.RespError { return response.NewError(err) }
	DefaultPageSize    = 10
)

type ActiveRecord interface{}

func DBClient() *gorm.DB {

	client, err := db.GetMySQL(conf.Config.MySQL["panda"])
	if err != nil {
		log.Fatal("db connection fatal ", err)
	}
	if conf.Config.Debug {
		client.LogMode(true)
	}
	DefaultDBClient = client
	return client
}

func RedisClient() *redis.Client {
	DefaultRedisClient = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		MinIdleConns: 5,
		PoolSize:     200,
	})

	_, err := DefaultRedisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	return DefaultRedisClient
}

func TxRecovery(tx *gorm.DB) {
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("recovery err:", r)
		}
	}()
}

func Lock(key string, expire time.Duration) (locked bool) {
	locked, _ = DefaultRedisClient.SetNX(key, 1, expire).Result()
	return
}

func UnLock(key string) {
	DefaultRedisClient.Del(key).Result()
}
