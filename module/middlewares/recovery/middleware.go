package recovery

import (
	"Panda/common/log"
	"Panda/common/response"
	"fmt"
	"github.com/kataras/iris/v12/context"
	"runtime"
	"strconv"
)

func getRequestLogs(ctx context.Context) string {
	var status, ip, method, path string
	status = strconv.Itoa(ctx.GetStatusCode())
	path = ctx.Path()
	method = ctx.Method()
	ip = ctx.RemoteAddr()
	// the date should be logged by iris' Logger, so we skip them
	return fmt.Sprintf("%v %s %s %s", status, path, method, ip)
}

func New() context.Handler {
	return func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ctx.IsStopped() {
					return
				}

				var stacktrace string
				for i := 1; ; i++ {
					_, f, l, got := runtime.Caller(i)
					if !got {
						break
					}
					stacktrace += fmt.Sprintf("%s:%d\n", f, l)
				}

				// when stack finishes
				logMessage := fmt.Sprintf("\nuser:%v,device:%v,version:%v\n", ctx.Values().Get("uid"), ctx.GetHeader("DEVICEID"),ctx.URLParamIntDefault("version_code", 0))
				logMessage += fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName())
				logMessage += fmt.Sprintf("At Request: %s\n", getRequestLogs(ctx))
				logMessage += fmt.Sprintf("Trace: %s\n", err)
				logMessage += fmt.Sprintf("\n%s", stacktrace)
				ctx.Application().Logger().Warn(logMessage)
				log.Error(logMessage)
				if ctx.GetStatusCode() == 200 {
					ctx.StatusCode(500)
					response.New(ctx).Error(10003).Format()
				} else {
					response.New(ctx).Error(10003).Format()
					//response.New(ctx).Error(10003).SetMessage("服务器已升级完毕，请重新登录").Format()
				}
				ctx.StopExecution()
			}
		}()
		ctx.Next()
	}
}
