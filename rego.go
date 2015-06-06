package rego

import (
	"net/http"
	"net/url"
)

type Server struct {
	*router
}

type Context struct {
	Params url.Values
}

type Result interface{}

type HandlerFunc func(c *Context) Result

func New() *Server {
	r := &router{mux: http.NewServeMux(), dispatchers: make(map[string]*dispatcher)}
	s := &Server{router: r}
	return s
}

func (s *Server) Run(addr string) {
	s.setHandler()
	if err := http.ListenAndServe(addr, s); err != nil {
		panic(err)
	}
}
