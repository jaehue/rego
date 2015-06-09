package rego

import (
	"net/http"
	"net/url"
)

type Server struct {
	*router
	middlewares []Middleware
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

func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) Run(addr string) {
	s.router.setHandler()
	if err := http.ListenAndServe(addr, s); err != nil {
		panic(err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	final := s.router.ServeHTTP

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		final = s.middlewares[i](final)
	}

	final(w, r)
}
