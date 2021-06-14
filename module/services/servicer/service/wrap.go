package service

import (
	"Panda/module/services/member/cache"
	"Panda/module/services/servicer"
	"github.com/go-redis/redis"
	"github.com/kataras/iris/v12"
	gocache "github.com/patrickmn/go-cache"
	"time"

	"Panda/module/models/member"
	members "Panda/module/services/member"
)

type ServicesWrap struct {
	userCache *cache.UserCache
	goCache   *gocache.Cache
	ctx       iris.Context
	user      member.UserStore
}

var _ servicer.Services = (*ServicesWrap)(nil)

func NewServicesWrap(r *redis.Client, uc *cache.UserCache) servicer.Services {
	s := new(ServicesWrap)
	s.userCache = uc
	s.goCache = gocache.New(30*time.Minute, 10*time.Minute)
	s.user = members.NewUser(s, r)
	return s
}

func (s *ServicesWrap) User() member.UserStore {
	return s.user
}
