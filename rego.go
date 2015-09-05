package rego

import "net/http"

type Server struct {
	*router
	middlewares []Middleware
	handlerFunc HandlerFunc
}

type App struct {
	Params map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type HandlerFunc func(*App)

func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) Run(addr string) {
	s.handlerFunc = s.router.handler()

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.handlerFunc = s.middlewares[i](s.handlerFunc)
	}

	if err := http.ListenAndServe(addr, s); err != nil {
		panic(err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a := &App{Params: make(map[string]interface{}), ResponseWriter: w, Request: r}
	for k, v := range r.URL.Query() {
		a.Params[k] = v[0]
	}

	s.handlerFunc(a)
}

func New() *Server {
	r := &router{dispatchers: make(map[string]*dispatcher)}
	s := &Server{router: r}
	s.middlewares = []Middleware{logHandler, recoverHandler, staticHandler, parseFormHandler, parseJsonBodyHandler}
	return s
}
