package member

import (
	"Panda/common/response"
	"context"
)

type UserStore interface {
	GetUser() (*UserInfo, response.RespError)
	CronTask()
	MqConsumer(ctx context.Context, maxWork int) <-chan struct{}
}
