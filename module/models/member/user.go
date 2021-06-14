package member

import (
	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
	"time"
)

type User struct {
	UserInfo
	redis *redis.Client `gorm:"-"`
	Cache *cache.Cache  `gorm:"-"`
}

type UserInfo struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  int    `json:"sex"`

}

func NewUser(r *redis.Client) *User {
	u := &User{redis: r, Cache: cache.New(30*time.Minute, 10*time.Minute)}
	return u
}
