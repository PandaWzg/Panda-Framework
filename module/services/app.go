package services

import (
	"sync"
)

var (
	Conf *appConf
)

type appConf struct {
	sync sync.RWMutex
}

func NewAppConf() *appConf {
	a := &appConf{
	}
	Conf = a
	return a
}

func (a *appConf) LoadAllConfig() {
	a.LoadConfig()
}

func (a *appConf) LoadConfig() {
	a.sync.Lock()
	//todo load config to cache
	a.sync.Unlock()
}
