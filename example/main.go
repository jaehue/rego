package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/jaehue/rego"
)

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	s := rego.New()
	s.HandleFunc("GET", "/", Index)
	s.HandleFunc("GET", "/users", Users)
	s.HandleFunc("GET", "/users/:id", func(c *rego.Context) {
		c.RenderJson(c.Params)
	})
	s.HandleFunc("POST", "/users", PostUser)

	s.HandleFunc("GET", "/login", func(c *rego.Context) {
		c.RenderTemplate("/public/login.html")
	})
	s.HandleFunc("POST", "/login", func(c *rego.Context) {
		if c.Params["username"] == "test" && c.Params["password"] == "password" {
			http.SetCookie(c.ResponseWriter, &http.Cookie{
				Name:  "X_AUTH",
				Value: Sign("verified"),
				Path:  "/",
			})
			c.Redirect("/public/index.html")
		}
		c.RenderTemplate("/public/login.html")

	})

	s.Use(AuthHandler)

	s.Run(":8082")
}

func Index(c *rego.Context) {
	c.RenderJson("Welcome rego")
}

func PostUser(c *rego.Context) {
	if u, ok := c.Params["user"]; ok {
		log.Println(u)
	}
}

func Users(c *rego.Context) {
	users := []User{User{1, "John", "john@mail.com"}, User{2, "Bob", "bob@mail.com"}, User{3, "Mark", "mark@mail.com"}}
	c.RenderJson(users)
}

func AuthHandler(next rego.HandlerFunc) rego.HandlerFunc {
	ignore := []string{"/login"}
	return func(c *rego.Context) {
		for _, s := range ignore {
			if strings.HasPrefix(c.Request.URL.Path, s) {
				next(c)
				return
			}
		}

		if v, err := c.Cookie("X_AUTH"); err == http.ErrNoCookie {
			// not authenticated
			c.Redirect("/login")
			return
		} else if err != nil {
			c.RenderErr(http.StatusInternalServerError, err)
			return
		} else if Verify("verified", v.Value) {
			// success
			next(c)
			return
		}
		c.Redirect("/login")
	}
}

func Sign(message string) string {
	secretKey := []byte("golang-book-secret-key2")
	if len(secretKey) == 0 {
		return ""
	}
	mac := hmac.New(sha1.New, secretKey)
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify returns true if the given signature is correct for the given message.
// e.g. it matches what we generate with Sign()
func Verify(message, sig string) bool {
	return hmac.Equal([]byte(sig), []byte(Sign(message)))
}
