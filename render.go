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

func (a *App) RenderTemplate(path string) {
	t, ok := loader.templates[path]
	if !ok {
		loader.once.Do(func() {
			t = template.Must(template.ParseFiles(filepath.Join(".", path)))
		})
		loader.templates[path] = t
	}
	t.Execute(a.ResponseWriter, nil)
}

func (a *App) RenderJson(v interface{}) {
	a.ResponseWriter.WriteHeader(http.StatusOK)
	a.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(a.ResponseWriter).Encode(v); err != nil {
		a.RenderErr(http.StatusInternalServerError, err)
	}
}

func (a *App) RenderXml(v interface{}) {
	a.ResponseWriter.WriteHeader(http.StatusOK)
	a.ResponseWriter.Header().Set("Content-Type", "application/xml; charset=utf-8")

	if err := xml.NewEncoder(a.ResponseWriter).Encode(v); err != nil {
		a.RenderErr(http.StatusInternalServerError, err)
	}
}

func (a *App) RenderErr(code int, err error) {
	if err != nil {
		if code > 0 {
			http.Error(a.ResponseWriter, http.StatusText(code), code)
		} else {
			defaultErr := http.StatusInternalServerError
			http.Error(a.ResponseWriter, http.StatusText(defaultErr), defaultErr)
		}
	}
}

func (a *App) Redirect(url string) {
	http.Redirect(a.ResponseWriter, a.Request, url, http.StatusMovedPermanently)
}

func (a *App) SetCookie(k, v string) {
	http.SetCookie(a.ResponseWriter, &http.Cookie{
		Name:  k,
		Value: v,
		Path:  "/",
	})
}

func (a *App) Cookie(k string) (*http.Cookie, error) {
	return a.Request.Cookie(k)
}
