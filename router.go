package rego

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	mux *http.ServeMux
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
	if ok && dispatcher.dispatch(w, req){
		return
	}

	for path, handler := range r.staticFileHandler {
		if strings.HasPrefix(req.URL.Path, path) {
			handler(w, req)
			return
		}

	}

	http.NotFound(w, req)
}

func (d *dispatcher) dispatch(w http.ResponseWriter, req *http.Request) bool {
	fn, ok := d.handles[req.URL.Path]
	if !ok {
		return false
	}

	a := &App{Params: make(map[string]interface{}), ResponseWriter: w, Request: req}
	if params, ok := ctx.GetAll(req); ok {
		for k, v := range params {
			a.Params[k] = v
		}
	}

	result := fn(a)
	if renderer, ok := result.(renderer); ok {
		renderer.render(w, req)
		return false
	}
	fmt.Fprint(w, result)
	return true
}
