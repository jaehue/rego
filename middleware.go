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

type Middleware func(next HandlerFunc) HandlerFunc

func logHandler(next HandlerFunc) HandlerFunc {
	return func(a *App) {
		t := time.Now()
		next(a)
		log.Printf("[%s] %q %v\n", a.Request.Method, a.Request.URL.String(), time.Now().Sub(t))
	}
}

func recoverHandler(next HandlerFunc) HandlerFunc {
	return func(a *App) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(a.ResponseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next(a)
	}
}

func staticHandler(next HandlerFunc) HandlerFunc {
	var (
		dir       = http.Dir(".")
		indexFile = "index.html"
	)
	return func(a *App) {
		r := a.Request
		w := a.ResponseWriter

		if r.Method != "GET" && r.Method != "HEAD" {
			next(a)
			return
		}

		file := r.URL.Path
		f, err := dir.Open(file)
		if err != nil {
			// discard the error?
			next(a)
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			next(a)
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
				next(a)
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				next(a)
				return
			}
		}
		http.ServeContent(w, r, file, fi.ModTime(), f)
	}
}

func parseJsonBodyHandler(next HandlerFunc) HandlerFunc {
	return func(a *App) {
		var m map[string]interface{}
		if json.NewDecoder(a.Request.Body).Decode(&m); len(m) > 0 {
			for k, v := range m {
				a.Params[k] = v
			}
		}
		next(a)
	}
}

func parseFormHandler(next HandlerFunc) HandlerFunc {
	return func(a *App) {
		a.Request.ParseForm()

		for k, v := range a.Request.PostForm {
			if len(v) > 0 {
				a.Params[k] = v[0]
			}
		}
		next(a)
	}
}
