package member

import (
	"Panda/common/helper/pool"
	"Panda/common/log"
	"Panda/common/queue/redis"
	"Panda/common/response"
	"Panda/module/models/member"
	"Panda/module/services"
	"Panda/module/services/servicer"
	"context"
	"encoding/json"
	"fmt"
	goredis "github.com/go-redis/redis"
	"time"
)

type syncFriend struct {
	uid   int
	level int8
}

type UserService struct {
	services servicer.Services
	user     *member.User
	redis    *goredis.Client
}

func NewUser(s servicer.Services, r *goredis.Client) member.UserStore {
	u := &UserService{
		services: s,
		redis:    r,
		user:     member.NewUser(r),
	}
	return u
}

func (u *UserService) GetUser() (*member.UserInfo, response.RespError) {
	var userInfo member.UserInfo
	userInfo.Name = "Panda"
	return &userInfo, nil
}

func (u *UserService) CronTask() {
	log.Info("cron task demo....")
	fmt.Println("cron task.....")
}

func (u *UserService) MqConsumer(ctx context.Context, maxWork int) <-chan struct{} {
	finish := make(chan struct{})
	queue := redis.NewQueue(services.DefaultRedisClient).Consumer("user_queue", ctx)
	go func() {
		defer func() {
			finish <- struct{}{}
			close(finish)
		}()
		pl := pool.New(maxWork) //max work
		for {
			select {
			case <-ctx.Done():
				for {
					select {
					case msg, ok := <-queue.GetMessage():
						if ok {
							pl.AddOne()
							u.execQueue(pl, msg.Value)
						}
					case <-time.After(2 * time.Second):
						goto LOOP
					}
				}
			case msg, ok := <-queue.GetMessage():
				if ok {
					pl.AddOne()
					u.execQueue(pl, msg.Value)
				}
			}
		}
	LOOP:
		pl.Wait()
	}()
	return finish
}

func (a *UserService) execQueue(pl *pool.Pool, value []byte) {
	var v member.User
	err := json.Unmarshal(value, &v)
	if err != nil {
		log.Error("message unmarshal error ", err)
	} else {
		go func(val member.User) {
			defer pl.DelOne()
			err = a.Receive(&val)
			if err == nil {
				//todo
			} else {
				//todo
			}
		}(v)
	}
}

func (a *UserService) Receive(params *member.User) (err error) {
	//todo 处理消息
	return
}
