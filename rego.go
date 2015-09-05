package rego

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

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

type templateLoader struct {
	once      sync.Once
	templates map[string]*template.Template
}

var loader = templateLoader{templates: make(map[string]*template.Template)}

func (a *App) SetCookie(k, v string) {
	http.SetCookie(a.ResponseWriter, &http.Cookie{
		Name:  k,
		Value: v,
		Path:  "/",
	})
}

func (a *App) RenderTemplate(path string) {
	t, ok := loader.templates[path]
	if !ok {
		loader.once.Do(func() {
			t = template.Must(template.ParseFiles(filepath.Join(".", path)))
		})
		loader.templates[path] = t
	}
	t.Execute(a.ResponseWriter, nil)
}

func (a *App) Redirect(url string) {
	http.Redirect(a.ResponseWriter, a.Request, url, http.StatusMovedPermanently)
}

func (a *App) RenderJson(v interface{}) {
	a.ResponseWriter.WriteHeader(http.StatusOK)
	a.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(a.ResponseWriter).Encode(v); err != nil {
		a.RenderErr(http.StatusInternalServerError, err)
	}
}

func (a *App) RenderXml(v interface{}) {
	a.ResponseWriter.WriteHeader(http.StatusOK)
	a.ResponseWriter.Header().Set("Content-Type", "application/xml; charset=utf-8")

	if err := xml.NewEncoder(a.ResponseWriter).Encode(v); err != nil {
		a.RenderErr(http.StatusInternalServerError, err)
	}
}

func (a *App) RenderErr(code int, err error) {
	if err != nil {
		if code > 0 {
			http.Error(a.ResponseWriter, http.StatusText(code), code)
		} else {
			defaultErr := http.StatusInternalServerError
			http.Error(a.ResponseWriter, http.StatusText(defaultErr), defaultErr)
		}
	}
}

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
