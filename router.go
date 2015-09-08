package rego

import (
	"net/http"
	"strings"
)

type router struct {
	// mux *http.ServeMux
	// key: HTTP Method (GET|POST|PUT|PATCH|DELETE)
	// value: HTTP Method별로 처리할 dispatcher
	dispatchers map[string]*dispatcher
}

type dispatcher struct {
	// key: URL Pattern
	// value:  URL Pattern별로 실행 할 HandlerFunc
	handles map[string]HandlerFunc
}

func (r *router) HandleFunc(method, pattern string, h HandlerFunc) {
	d, ok := r.dispatchers[method]
	if !ok {
		d = &dispatcher{make(map[string]HandlerFunc)}
		r.dispatchers[method] = d
	}
	d.handles[pattern] = h
}

func (r *router) handler() HandlerFunc {
	return func(c *Context) {
		// 요청 Method에 해당하는 HandlerFunc를 호출
		dispatcher, ok := r.dispatchers[c.Request.Method]
		if ok && dispatcher.dispatch(c) {
			return
		}

		http.NotFound(c.ResponseWriter, c.Request)
	}
}

func (d *dispatcher) dispatch(c *Context) bool {
	fn, params, found := d.lookup(c.Request.URL.Path)
	if !found {
		return false
	}

	for k, v := range params {
		c.Params[k] = v
	}

	fn(c)
	return true
}

func (d *dispatcher) lookup(path string) (HandlerFunc, map[string]string, bool) {
	for pattern, handler := range d.handles {
		if matched, params := match(pattern, path); matched {
			return handler, params, true
		}
	}
	return nil, nil, false
}

func match(pattern, path string) (matched bool, params map[string]string) {
	if pattern == path {
		return true, map[string]string{}
	}
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	if len(patterns) != len(paths) {
		return
	}
	params = make(map[string]string)
	for i := 0; i < len(patterns); i++ {
		switch {
		case len(patterns[i]) == 0:
		case patterns[i] == paths[i]:
			matched = true
		case patterns[i][0] == ':':
			params[patterns[i][1:]] = paths[i]
			matched = true
		default:
			matched = false
			break
		}
	}
	return
}
