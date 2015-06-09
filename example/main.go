package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jaehue/rego"
)

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	s := rego.New()
	s.Get("/", Index)
	s.Get("/users", Users)
	s.Use(logHandler)
	s.Run(":8082")
}

func Index(c *rego.Context) rego.Result {
	return "Welcome rego"
}

func Users(c *rego.Context) rego.Result {
	users := []User{User{1, "John", "john@mail.com"}, User{2, "Bob", "bob@mail.com"}, User{3, "Mark", "mark@mail.com"}}
	return c.RenderJson(users)
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), time.Now().Sub(t))
	})
}
