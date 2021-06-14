package api

import (
	controllers2 "Panda/api/controllers"
	"Panda/common/cron"
	"Panda/common/log"
	"Panda/common/response"
	"Panda/conf"
	"Panda/module/middlewares/access"
	"Panda/module/middlewares/recovery"
	"Panda/module/services"
	cache2 "Panda/module/services/member/cache"
	"Panda/module/services/servicer/service"
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

type api struct {
	config *conf.Cfg
	cron   bool
}

func New(cfg *conf.Cfg) *api {
	return &api{config: cfg}
}

func (r *api) StartCron() *api {
	r.cron = true
	return r
}

func (r *api) Start(stdCtx context.Context) error {
	app := iris.New()
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		log.Infof("request path:%v not found, method:%v params:%v", ctx.Path(), ctx.Method(), ctx.URLParams())
		response.New(ctx).Error(10000).Format()
	})

	//初始化客户端连接
	services.DBClient()
	services.RedisClient()

	//init load config to cache
	appConf := services.NewAppConf()
	appConf.LoadAllConfig()

	//供全局任意地方调
	uc := cache2.NewUserCache(services.DefaultRedisClient)
	serviceWrap := service.NewServicesWrap(services.DefaultRedisClient, uc)

	mvcApp := mvc.New(app)
	api := mvcApp.Party("/")
	api.Register(serviceWrap)

	//加入中间件
	api.Router.Use(recovery.New()) //捕获异常
	api.Router.Use(access.New())   //安全校验

	api.Party("user").Handle(new(controllers2.UserController))

	//服务启动开始消费mq
	ctx, cancel := context.WithCancel(context.Background())
	finish := serviceWrap.User().MqConsumer(ctx, 500)

	if r.cron {
		c := cron.New()
		c.Start()
		defer c.Stop()
		_ = c.AddFunc("0 */30 * * * ?", serviceWrap.User().CronTask)

	}
	iris.RegisterOnInterrupt(func() {
		fmt.Println("start shutdown...")
		cancel()
		<-finish
		_ = app.Shutdown(ctx)
	})

	return app.Run(iris.Addr(":"+strconv.Itoa(r.config.Frontend.Port)), iris.WithoutInterruptHandler, iris.WithoutServerError(iris.ErrServerClosed), iris.WithRemoteAddrHeader("X-Real-IP"))

}
