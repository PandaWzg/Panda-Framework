package signature

import (
	"Panda/common/response"
	"github.com/kataras/iris/v12/context"
)

func New() context.Handler {
	return func(ctx context.Context) {
		params := make(map[string]string, 0)
		params = ctx.URLParams()

		if ctx.Method() == "POST" {
			formValues := ctx.FormValues()
			for k, v := range formValues {
				params[k] = v[0]
			}

			body, err := context.GetBody(ctx.Request(), true)
			//body, err := ctx.GetBody()
			if err != nil {
				response.New(ctx).Error(20001).Format()
				return
			}
			if len(body) > 0 {
				//todo 签名校验

			}

		} else {
			params = ctx.URLParams()
		}

		ctx.Next()
	}
}
