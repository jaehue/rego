package rego

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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
	err error
}

func (r jsonResult) render(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(r.v); err != nil {
		errResult{err}.render(w, req)
	}
}

func (r xmlResult) render(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	if err := xml.NewEncoder(w).Encode(r.v); err != nil {
		errResult{err}.render(w, req)
	}
}

func (r errResult) render(w http.ResponseWriter, req *http.Request) {
	if r.err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, r.err)
	}
}
