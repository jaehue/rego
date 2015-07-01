package rego

import "net/http"

type Server struct {
	*router
	middlewares []Middleware
}

type Context struct {
	Params map[string]interface{}
}

type Result interface{}

type HandlerFunc func(c *Context) Result

func (c *Context) RenderJson(v interface{}) Result {
	return jsonResult{v}
}

func (c *Context) RenderXml(v interface{}) Result {
	return xmlResult{v}
}

func (c *Context) RenderErr(code int, err error) Result {
	return errResult{code, err}
}

func New() *Server {
	r := &router{mux: http.NewServeMux(), dispatchers: make(map[string]*dispatcher)}
	s := &Server{router: r}
	s.middlewares = []Middleware{logHandler, AuthHandler, bodyParserHandler}
	return s
}

func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) Run(addr string) {
	s.router.setHandler()

	var final http.Handler = s.router

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		final = s.middlewares[i](final)
	}

	if err := http.ListenAndServe(addr, final); err != nil {
		panic(err)
	}
}
