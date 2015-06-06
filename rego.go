package rego

import (
	"fmt"
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

func (s *Server) Get(path string, h HandlerFunc) {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		c := &Context{}
		result := h(c)
		fmt.Fprint(w, result)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Run(addr string) {
	if err := http.ListenAndServe(addr, s); err != nil {
		panic(err)
	}
}

func New() *Server {
	mux := http.NewServeMux()
	s := &Server{mux: mux}
	return s
}
