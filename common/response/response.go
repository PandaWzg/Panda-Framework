package response

import (
	"github.com/kataras/iris/v12"
	"reflect"
)

type Result struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    interface{}            `json:"data"`
	Error   interface{}            `json:"-"`
	Format  string                 `json:"-"`
	DataMap map[string]interface{} `json:"-"`
}

type RespError interface {
	GetCode() int
	GetError() interface{}
}

type Response struct {
	ctx    iris.Context
	Result Result
}

const (
	jsonType      = "json"
	xmlType       = "xml"
	defaultFormat = "json"
	itemField     = "items"
	countField    = "count"
)

func NewCode(code int) RespError {
	r := &Response{}
	r.Result.Code = code
	return r
}

func NewError(err interface{}) RespError {
	r := &Response{}
	r.Result.Error = err
	return r
}

func NotFound() RespError {
	r := &Response{}
	r.Result.Code = 0
	return r
}

func New(ctx iris.Context) *Response {
	r := &Response{
		ctx: ctx,
	}
	r.Result.Code = 200
	r.Result.Format = defaultFormat
	r.Result.DataMap = make(map[string]interface{}, 0)
	return r
}

func (r *Response) GetCode() int {
	return r.Result.Code
}

func (r *Response) GetError() interface{} {
	return r.Result.Error
}

func (r *Response) SetJson() *Response {
	r.Result.Format = jsonType
	return r
}

func (r *Response) SetXml() *Response {
	r.Result.Format = xmlType
	return r
}

func (r *Response) Error(code int) *Response {
	r.Result.Code = code
	r.Result.Data = nil
	r.Result.Message = GetError(code)
	return r
}

func (r *Response) Data(data interface{}) *Response {
	if data != nil {
		r.Result.Data = data
	} else {
		r.Result.Data = nil
	}
	return r
}

func interfaceLen(args interface{}) int {
	val := reflect.ValueOf(args)
	return val.Len()
}

func (r *Response) Items(items interface{}) *Response {
	if interfaceLen(items) > 0 {
		r.Result.DataMap[itemField] = items
	} else {
		r.Result.DataMap[itemField] = []string{}
	}
	r.Result.Data = r.Result.DataMap

	return r
}

func (r *Response) Count(count int64) *Response {
	r.Result.DataMap[countField] = count
	r.Result.Data = r.Result.DataMap
	return r
}

func (r *Response) SetMessage(msg string) *Response {
	r.Result.Message = msg
	return r
}

func (r *Response) Format() {
	if r.Result.Message == "" {
		r.SetMessage("success")
	}

	if r.Result.Format == jsonType {
		if _, err := r.ctx.JSON(r.Result); err != nil {
			panic(err)
		}
	} else {
		if _, err := r.ctx.XML(r.Result); err != nil {
			panic(err)
		}
	}
}
