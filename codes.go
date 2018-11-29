package httpserver

const (
	OK         = 0
	ErrBinding = 1
	ErrAck     = 2
	ErrServer  = 3
)

var (
	builtinErr = map[int]string{
		OK:         "请求成功",
		ErrBinding: "数据绑定错误",
		ErrAck:     "服务未返回数据",
	}
)

func errorMessage(code int) string {
	return builtinErr[code]
}
