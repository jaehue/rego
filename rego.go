package rego

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

type Server struct {
	*router
	middlewares []Middleware
}

type App struct {
	Params map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func (a App) SetCookie(k, v string) {
	http.SetCookie(a.ResponseWriter, &http.Cookie{
		Name:  k,
		Value: v,
		Path:  "/",
	})
}

type Result interface{}

type HandlerFunc func(*App) Result

type templateLoader struct {
	once      sync.Once
	templates map[string]*template.Template
}

var loader = templateLoader{templates: make(map[string]*template.Template)}

func (a *App) RenderTemplate(path string) Result {
	t, ok := loader.templates[path]
	if !ok {
		loader.once.Do(func() {
			t = template.Must(template.ParseFiles(filepath.Join(".", path)))
		})
		loader.templates[path] = t
	}
	return templateResult{t}
}

func (a *App) Redirect(url string) Result {
	a.ResponseWriter.Header().Set("Location", url)
	a.ResponseWriter.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

func (a *App) RenderJson(v interface{}) Result {
	return jsonResult{v}
}

func (a *App) RenderXml(v interface{}) Result {
	return xmlResult{v}
}

func (a *App) RenderErr(code int, err error) Result {
	return errResult{code, err}
}

func New() *Server {
	r := &router{dispatchers: make(map[string]*dispatcher)}
	s := &Server{router: r}
	s.middlewares = []Middleware{logHandler, parseFormHandler, parseJsonBodyHandler}
	return s
}

func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) Run(addr string) {
	var final http.Handler = s.router

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		final = s.middlewares[i](final)
	}

	if err := http.ListenAndServe(addr, final); err != nil {
		panic(err)
	}
}
