package controllers

import (
	"Panda/common/log"
	"Panda/common/response"
	"Panda/conf"
	"Panda/module/models/member"
	"github.com/kataras/iris/v12"
	"strings"
)

type controller struct {
	Ctx    iris.Context
	Config *conf.Cfg
	User   *member.User
}

func (c *controller) successCode(code int) {
	response.New(c.Ctx).Data(nil).SetMessage(response.GetError(code)).Format()
}

func (c *controller) success(result ...interface{}) {
	if len(result) == 1 {
		response.New(c.Ctx).Data(result[0]).Format()
	} else if len(result) > 1 {
		resp := response.New(c.Ctx).Data(result[0])
		if code, ok := result[1].(int); ok {
			resp.SetMessage(response.GetError(code))
		}
		resp.Format()
	} else {
		response.New(c.Ctx).Data(nil).Format()
	}
}


func (c *controller) errorCode(code int) {
	response.New(c.Ctx).Error(code).Format()
}

func (c *controller) parseParams(v interface{}) {
	if err := c.Ctx.ReadJSON(&v); err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		panic(err)
	}
}

// errorResp 使用方式:
//
// 1. errorCode(10001)			参数10001为错误代码
// 2. errorCode(10001, err)		参数10001为错误代码; err类型为error,string,可多个
// 3. errorCode(respErr)		参数respErr类型为response.RespError
// 4. errorCode(respErr, err)	参数respErr类型为response.RespError; err类型为error,string,可多个
//
// 当全局debug开启时，第二个及之后的错误信息会替换返回中的message，否则输出已配置的code错误信息，并把第二个及之后的错误信息输出到log
// 等于0的code会被替换为10001
func (c *controller) errorResp(respCode interface{}, err ...interface{}) {
	code := 0
	messages := make([]string, 0)

	switch respCode.(type) {
	case int:
		code = respCode.(int)
	case response.RespError:
		code = respCode.(response.RespError).GetCode()
		if code == 0 {
			err := respCode.(response.RespError).GetError()
			switch err.(type) {
			case error:
				messages = append(messages, err.(error).Error())
			case string:
				messages = append(messages, err.(string))
			}
		}
	default:
		panic("unsupported respCode type")
	}

	for _, e := range err {
		switch e.(type) {
		case string:
			messages = append(messages, e.(string))
		case error:
			messages = append(messages, e.(error).Error())
		default:
			panic("unsupported err type")
		}
	}

	if code == 0 {
		code = 10001
	}

	resp := response.New(c.Ctx).Error(code)
	if len(messages) > 0 {
		resp.SetMessage(strings.Join(messages, "; "))
	}
	log.Error(strings.Join(messages, "; "))
	resp.Format()
}
