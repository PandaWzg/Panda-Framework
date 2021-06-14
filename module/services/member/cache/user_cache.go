package cache

import (
	"github.com/go-redis/redis"
	"sync"
)

type UserCache struct {
	redis *redis.Client
	sync  sync.RWMutex
}

type Stat struct {
	Total    int64
	DayCount int64
}

func NewUserCache(r *redis.Client) *UserCache {
	return &UserCache{redis: r}
}
//todo 频控相关code