package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	e *gin.Engine
}

func NewHttpServer() *HttpServer {

	var (
		s = &HttpServer{
			e: gin.Default(),
		}
	)

	return s
}

func (m *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.e.ServeHTTP(w, r)
}

func (m *HttpServer) Handle(httpMethod string, relativePath string, fn interface{}) {

	var (
		meta = &Meta{
			Func: fn,
		}
	)

	meta.check(relativePath)

	m.e.Handle(httpMethod, "apis/"+relativePath, meta.handler)

}

var (
	defaultServer = NewHttpServer()
)

func GET(relativePath string, fn interface{}) {
	defaultServer.Handle("GET", relativePath, fn)
}

func POST(relativePath string, fn interface{}) {
	defaultServer.Handle("POST", relativePath, fn)
}

func ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, defaultServer)
}
