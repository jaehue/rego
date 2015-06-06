package rego

import (
	"fmt"
	"net/http"
)

func (s *Server) Get(path string, h HandlerFunc) {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		c := &Context{}
		result := h(c)
		fmt.Fprint(w, result)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
