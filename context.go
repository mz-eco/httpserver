package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mz-eco/types"
)

type contextGetter interface {
	context() *Context
}

type Context struct {
	ctx *gin.Context
	ack interface{}
}

var (
	TypeContext = types.ElemType((*Context)(nil))
)

func (m *Context) Error(code int, err error) {

	fmt.Println(err.Error())

	m.ctx.JSON(
		http.StatusOK,
		map[string]interface{}{
			"Code":    code,
			"Message": errorMessage(code),
		},
	)
}

func (m *Context) Message(format string, v ...interface{}) {

	m.ctx.JSON(
		http.StatusOK,
		map[string]interface{}{
			"Code":    ErrServer,
			"Message": fmt.Sprintf(format, v...),
		},
	)
}

func (m *Context) Done(v interface{}) {

	if v == nil {
		m.ack = make(map[string]interface{})
	} else {
		m.ack = v
	}

}

func (m *Context) init(ctx *gin.Context) {
	m.ctx = ctx
}

func (m *Context) context() *Context {
	return m
}
