package httpserver

import (
	"fmt"
	"net/http"

	WS "github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"github.com/mz-eco/types"
)

var (
	grader = WS.Upgrader{}
)

type contextGetter interface {
	context() *Context
}

type Context struct {
	ctx *gin.Context
	ack interface{}
}

func (m *Context) Socket() (*WS.Conn, error) {
	return grader.Upgrade(m.ctx.Writer, m.ctx.Request, nil)
}

type ackErrorType struct {
}

var (
	ackError = ackErrorType{}
)

var (
	TypeContext = types.ElemType((*Context)(nil))
)

func (m *Context) Error(code int, err error) {
	m.ctx.JSON(
		http.StatusOK,
		map[string]interface{}{
			"Code":    code,
			"Message": errorMessage(code),
		},
	)

	m.ack = ackError
}

func (m *Context) Message(format string, v ...interface{}) {

	m.ctx.JSON(
		http.StatusOK,
		map[string]interface{}{
			"Code":    ErrServer,
			"Message": fmt.Sprintf(format, v...),
		},
	)

	m.ack = ackError
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
