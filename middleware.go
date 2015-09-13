package rego

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type Middleware func(next HandlerFunc) HandlerFunc

func logHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		next(c)
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), time.Now().Sub(t))
	}
}

func recoverHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(c.ResponseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next(c)
	}
}

func staticHandler(next HandlerFunc) HandlerFunc {
	var (
		dir       = http.Dir(".")
		indexFile = "index.html"
	)
	return func(c *Context) {
		r := c.Request
		w := c.ResponseWriter

		if r.Method != "GET" && r.Method != "HEAD" {
			next(c)
			return
		}

		file := r.URL.Path
		f, err := dir.Open(file)
		if err != nil {
			// discard the error?
			next(c)
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			next(c)
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
				next(c)
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				next(c)
				return
			}
		}
		http.ServeContent(w, r, file, fi.ModTime(), f)
	}
}

func parseJsonBodyHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		var m map[string]interface{}
		if json.NewDecoder(c.Request.Body).Decode(&m); len(m) > 0 {
			for k, v := range m {
				c.Params[k] = v
			}
		}
		next(c)
	}
}

func parseFormHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		c.Request.ParseForm()

		for k, v := range c.Request.PostForm {
			if len(v) > 0 {
				c.Params[k] = v[0]
			}
		}
		next(c)
	}
}
