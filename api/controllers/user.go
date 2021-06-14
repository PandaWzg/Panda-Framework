package controllers

import (
	"Panda/conf"
	"Panda/module/services/servicer"
	"github.com/kataras/iris/v12"
)

type InternalUserController struct {
	controller
	Ctx      iris.Context
	Config   *conf.Cfg
	Services servicer.Services
}

type UserController struct {
	controller
	Ctx      iris.Context
	Services servicer.Services
}

func (c *UserController) GetInfo() {
	if result, err := c.Services.User().GetUser(); err != nil {
		c.errorCode(err.GetCode())
	} else {
		c.success(result)
	}
}
