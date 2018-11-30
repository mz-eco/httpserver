package httpserver

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	e   *gin.Engine
	mux *http.ServeMux
}

func NewHttpServer() *HttpServer {

	//gin.SetMode(gin.ReleaseMode)

	var (
		s = &HttpServer{
			e:   gin.Default(),
			mux: &http.ServeMux{},
		}
	)

	s.e.Use(gin.Recovery())

	return s
}

func (m *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, "/api") {
		m.e.ServeHTTP(w, r)
	} else {
		m.mux.ServeHTTP(w, r)
	}

}

func (m *HttpServer) Any(relativePath string, fn interface{}) {

	var (
		meta = &Meta{
			Func: fn,
		}
	)

	meta.check(relativePath)

	m.e.Any("apis/"+relativePath, meta.handler)

}

func (m *HttpServer) Web(dir string) {

	dir, _ = filepath.Abs(dir)

	m.mux.Handle("/", http.FileServer(http.Dir(dir)))

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

func Web(dir string) {
	defaultServer.Web(dir)
}
func Any(relativePath string, fn interface{}) {
	defaultServer.Any(relativePath, fn)
}

func POST(relativePath string, fn interface{}) {
	defaultServer.Handle("POST", relativePath, fn)
}

func ListenAndServe(addr string) error {

	return http.ListenAndServe(addr, defaultServer)
}
