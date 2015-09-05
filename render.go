package rego

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

var loader = templateLoader{templates: make(map[string]*template.Template)}

type templateLoader struct {
	once      sync.Once
	templates map[string]*template.Template
}

func (c *Context) RenderTemplate(path string) {
	t, ok := loader.templates[path]
	if !ok {
		loader.once.Do(func() {
			t = template.Must(template.ParseFiles(filepath.Join(".", path)))
		})
		loader.templates[path] = t
	}
	t.Execute(c.ResponseWriter, nil)
}

func (c *Context) RenderJson(v interface{}) {
	c.ResponseWriter.WriteHeader(http.StatusOK)
	c.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

func (c *Context) RenderXml(v interface{}) {
	c.ResponseWriter.WriteHeader(http.StatusOK)
	c.ResponseWriter.Header().Set("Content-Type", "application/xml; charset=utf-8")

	if err := xml.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

func (c *Context) RenderErr(code int, err error) {
	if err != nil {
		if code > 0 {
			http.Error(c.ResponseWriter, http.StatusText(code), code)
		} else {
			defaultErr := http.StatusInternalServerError
			http.Error(c.ResponseWriter, http.StatusText(defaultErr), defaultErr)
		}
	}
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.ResponseWriter, c.Request, url, http.StatusMovedPermanently)
}

func (c *Context) SetCookie(k, v string) {
	http.SetCookie(c.ResponseWriter, &http.Cookie{
		Name:  k,
		Value: v,
		Path:  "/",
	})
}

func (c *Context) Cookie(k string) (*http.Cookie, error) {
	return c.Request.Cookie(k)
}
