package access

import (
	"Panda/common/log"
	"Panda/common/response"
	"Panda/conf"
	"github.com/kataras/iris/v12/context"
	"strings"
)

var defaultAllowIp = "10.*.*.*"
func New() context.Handler {
	return func(ctx context.Context) {
		remoteIp := ctx.RemoteAddr()
		params := ctx.URLParams()
		if _, ok := params["internal-no-sign"]; ok && conf.Config.Env == "prod" { //生产环境验证无sign访问来源
			if checkIpRules(remoteIp) == false {
				response.New(ctx).Error(40001).Format()
				return
			}
			log.Infof("no sign request info:%v, ip:%s, uri:%v", params, remoteIp, ctx.FullRequestURI())
		}
		ctx.Next()
	}
}

func checkIpRules(ip string) bool {
	ipRules := strings.Split(defaultAllowIp, ".")
	ips := strings.Split(ip, ".")
	if ips[0] == "::1" {
		return true
	}

	for k, v := range ipRules {
		if v != "*" && v != ips[k] {
			return false
		}
	}
	return true
}
