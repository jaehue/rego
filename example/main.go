package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jaehue/rego"
)

func main() {
	s := rego.New()
	s.Get("/", Index)
	s.Get("/hello", Hello)
	s.Use(logHandler)
	s.Run(":8082")
}

func Index(c *rego.Context) rego.Result {
	return "Welcome rego"
}

func Hello(c *rego.Context) rego.Result {
	return "Hello rego"
}

func logHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next(w, r)
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), time.Now().Sub(t))
	}
}
