package rego

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
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

func recoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func staticHandler(next http.Handler) http.Handler {
	var (
		dir       = http.Dir(".")
		indexFile = "index.html"
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "HEAD" {
			next.ServeHTTP(w, r)
			return
		}

		file := r.URL.Path
		f, err := dir.Open(file)
		if err != nil {
			// discard the error?
			next.ServeHTTP(w, r)
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// try to serve index file
		if fi.IsDir() {
			// redirect if missing trailing slash
			if !strings.HasSuffix(r.URL.Path, "/") {
				http.Redirect(w, r, r.URL.Path+"/", http.StatusFound)
				return
			}

			file = path.Join(file, indexFile)
			f, err = dir.Open(file)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.ServeContent(w, r, file, fi.ModTime(), f)
	})
}

func parseJsonBodyHandler(next http.Handler) http.Handler {
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

func parseFormHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for k, v := range r.PostForm {
			if len(v) > 0 {
				ctx.Set(r, k, v[0])
			}
		}
		next.ServeHTTP(w, r)

	})
}
