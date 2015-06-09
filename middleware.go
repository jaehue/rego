package rego

import "net/http"

type middlewareChain struct {
	middlewares []Middleware
}

type Middleware func(next http.Handler) http.Handler
