package rego

import (
	"fmt"
	"net/http"
)

type router struct {
	mux *http.ServeMux

	// key: URL Pattern
	// value: URL Pattern별로 처리할 dispatcher
	dispatchers map[string]*dispatcher
}

type dispatcher struct {
	// key: Method (GET|POST|PUT|PATCH|DELETE)
	// value:  method 별로 실행 할 HandlerFunc
	handles map[string]HandlerFunc
}

// 요청 Method에 해당하는 HandlerFunc를 호출
func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn, ok := d.handles[r.Method]
	if !ok {
		http.NotFound(w, r)
		return
	}

	c := &Context{}
	result := fn(c)
	fmt.Fprint(w, result)
}

func (r *router) Get(path string, h HandlerFunc) {
	r.register("GET", path, h)
}

func (r *router) Post(path string, h HandlerFunc) {
	r.register("POST", path, h)
}

func (r *router) register(method, pattern string, h HandlerFunc) {
	d, ok := r.dispatchers[pattern]
	if !ok {
		d = &dispatcher{make(map[string]HandlerFunc)}
		r.dispatchers[pattern] = d
	}
	d.handles[method] = h
}

func (r *router) setHandler() {
	for p, d := range r.dispatchers {
		r.mux.Handle(p, d)
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
