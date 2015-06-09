package rego

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type middlewareChain struct {
	middlewares []Middleware
}

type Middleware func(next http.Handler) http.Handler

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), time.Now().Sub(t))
	})
}

func bodyParserHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m map[string]interface{}
		if json.NewDecoder(r.Body).Decode(&m); len(m) > 0 {
			for k, v := range m {
				ctx.Set(r, k, v)
			}
		}
		next.ServeHTTP(w, r)
	})
}
