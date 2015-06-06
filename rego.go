package rego

import (
	"net/http"
	"net/url"
)

type Server struct {
	mux *http.ServeMux
}

type Context struct {
	Params url.Values
}

type Result interface{}

type HandlerFunc func(c *Context) Result

func New() *Server {
	mux := http.NewServeMux()
	s := &Server{mux: mux}
	return s
}

func (s *Server) Run(addr string) {
	if err := http.ListenAndServe(addr, s); err != nil {
		panic(err)
	}
}
