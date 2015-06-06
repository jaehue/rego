package rego

import (
	"net/http"
	"testing"
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header)           { return http.Header{} }
func (m *mockResponseWriter) Write(p []byte) (n int, err error) { return len(p), nil }
func (m *mockResponseWriter) WriteHeader(int)                   {}

func TestServer(t *testing.T) {
	s := New()

	routed := false
	s.Get("/", func(c *Context) Result {
		routed = true
		return nil
	})

	w := new(mockResponseWriter)
	r, _ := http.NewRequest("GET", "/", nil)
	s.mux.ServeHTTP(w, r)

	if !routed {
		t.Fatal("routing failed")
	}
}
