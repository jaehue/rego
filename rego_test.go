package rego

import "net/http"

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header)           { return http.Header{} }
func (m *mockResponseWriter) Write(p []byte) (n int, err error) { return len(p), nil }
func (m *mockResponseWriter) WriteHeader(int)                   {}
