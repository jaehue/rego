package main

import (
	"encoding/json"
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
	s.Post("/users", PostUser)
	s.Use(logHandler, bodyParserHandler)
	s.Run(":8082")
}

func Index(c *rego.Context) rego.Result {
	return "Welcome rego"
}

func PostUser(c *rego.Context) rego.Result {
	if u, ok := c.Params["user"]; ok {
		log.Println(u)
	}
	return nil
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

func bodyParserHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m map[string]interface{}
		if json.NewDecoder(r.Body).Decode(&m); len(m) > 0 {
			for k, v := range m {
				rego.Ctx.Set(r, k, v)
			}
		}
		next.ServeHTTP(w, r)
	})
}
