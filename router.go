package rego

import (
	"fmt"
	"net/http"
)

type router struct {
	// key: HTTP Method (GET|POST|PUT|PATCH|DELETE)
	// value: HTTP Method별로 처리할 dispatcher
	dispatchers       map[string]*dispatcher
	staticFileHandler map[string]http.HandlerFunc
}

type dispatcher struct {
	// key: URL Pattern
	// value:  URL Pattern별로 실행 할 HandlerFunc
	handles map[string]HandlerFunc
}

func (r *router) Get(path string, h HandlerFunc) {
	r.HandleFunc("GET", path, h)
}

func (r *router) Post(path string, h HandlerFunc) {
	r.HandleFunc("POST", path, h)
}

func (r *router) Static(path string) {
	if r.staticFileHandler == nil {
		r.staticFileHandler = make(map[string]http.HandlerFunc)
	}

	r.staticFileHandler[path] = func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, req.URL.Path[1:])
	}
}

func (r *router) HandleFunc(method, pattern string, h HandlerFunc) {
	d, ok := r.dispatchers[method]
	if !ok {
		d = &dispatcher{make(map[string]HandlerFunc)}
		r.dispatchers[method] = d
	}
	d.handles[pattern] = h
}

// 요청 Method에 해당하는 HandlerFunc를 호출
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	dispatcher, ok := r.dispatchers[req.Method]
	if !ok {
		http.NotFound(w, req)
		return
	}

	fn, ok := dispatcher.handles[req.URL.Path]
	if !ok {
		http.NotFound(w, req)
		return
	}

	c := &Context{Params: make(map[string]interface{}), ResponseWriter: w, Request: req}
	if params, ok := ctx.GetAll(req); ok {
		for k, v := range params {
			c.Params[k] = v
		}
	}

	result := fn(c)
	if renderer, ok := result.(renderer); ok {
		renderer.render(w, req)
		return
	}
	fmt.Fprint(w, result)
}
