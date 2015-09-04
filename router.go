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

// 요청 Method에 해당하는 HandlerFunc를 호출
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	dispatcher, ok := r.dispatchers[req.Method]
	if ok && dispatcher.dispatch(w, req) {
		return
	}

	http.NotFound(w, req)
}

func (d *dispatcher) dispatch(w http.ResponseWriter, req *http.Request) bool {
	fn, params, found := d.lookup(req.URL.Path)
	if !found {
		return false
	}

	a := NewApp(w, req, params)

	fn(a)
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
		if len(patterns[i]) == 0 {
			continue
		}
		if patterns[i] == paths[i] {
			matched = true
			continue
		}
		if patterns[i][0] == ':' {
			params[patterns[i][1:]] = paths[i]
			matched = true
			continue
		}
		matched = false
		break
	}
	return
}
