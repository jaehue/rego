package rego

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

func AuthHandler(next http.Handler) http.Handler {
	ignore := []string{"/login"}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, s := range ignore {
			if strings.HasPrefix(r.URL.Path, s) {
				next.ServeHTTP(w, r)
				return
			}
		}

		if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
			// not authenticated
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else if err != nil {
			// some other error
			panic(err.Error())
		} else {
			// success - call the next handler
			next.ServeHTTP(w, r)
		}
	})
}
