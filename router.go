package rego

import (
	"net/http"
	"strings"
)

type router struct {
	// key: HTTP Method (GET|POST|PUT|PATCH|DELETE)
	// value: URL Pattern별로 실행 할 HandlerFunc
	handlers map[string]map[string]HandlerFunc
}

func (r *router) HandleFunc(method, pattern string, h HandlerFunc) {
	d, ok := r.handlers[method]
	if !ok {
		d = make(map[string]HandlerFunc)
		r.handlers[method] = d
	}
	d[pattern] = h
}

func (r *router) handler() HandlerFunc {
	return func(c *Context) {
		for pattern, handler := range r.handlers[c.Request.Method] {
			if ok, params := match(pattern, c.Request.URL.Path); ok {
				if params != nil {
					for k, v := range params {
						c.Params[k] = v
					}
				}
				// 요청 url에 해당하는 handler 수행
				handler(c)
				return
			}
		}
		// 요청 url에 해당하는 handler를 찾지 못한 경우 NotFound 에러 처리
		http.NotFound(c.ResponseWriter, c.Request)
		return

	}
}

func match(pattern, path string) (bool, map[string]string) {
	if pattern == path {
		return true, nil
	}
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	if len(patterns) != len(paths) {
		return false, nil
	}
	params := make(map[string]string)
	for i := 0; i < len(patterns); i++ {
		switch {
		case len(patterns[i]) == 0:
		case patterns[i] == paths[i]:
		case patterns[i][0] == ':':
			params[patterns[i][1:]] = paths[i]
		default:
			return false, nil
		}
	}
	return true, params
}
