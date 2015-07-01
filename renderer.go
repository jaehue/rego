package rego

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
)

type renderer interface {
	render(w http.ResponseWriter, req *http.Request)
}

type jsonResult struct {
	v interface{}
}

type xmlResult struct {
	v interface{}
}

type errResult struct {
	code int
	err  error
}

type templateResult struct {
	*template.Template
}

func (r templateResult) render(w http.ResponseWriter, req *http.Request) {
	r.Execute(w, nil)
}

func (r jsonResult) render(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(r.v); err != nil {
		errResult{err: err}.render(w, req)
	}
}

func (r xmlResult) render(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")

	if err := xml.NewEncoder(w).Encode(r.v); err != nil {
		errResult{err: err}.render(w, req)
	}
}

func (r errResult) render(w http.ResponseWriter, req *http.Request) {
	if r.err != nil {
		if r.code > 0 {
			http.Error(w, http.StatusText(r.code), r.code)
		} else {
			defaultErr := http.StatusInternalServerError
			http.Error(w, http.StatusText(defaultErr), defaultErr)
		}
	}
}
