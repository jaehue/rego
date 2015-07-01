package main

import (
	"log"

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

	s.Static("/public/")
	s.Get("/login", func(c *rego.Context) rego.Result {
		return c.RenderTemplate("/public/login.html")
	})
	s.Post("/login", func(c *rego.Context) rego.Result {
		c.SetCookie("auth", "true")
		return c.Redirect("/public/index.html")
	})

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
