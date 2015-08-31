package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
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
	s.HandleFunc("POST", "/users", PostUser)

	s.Static("/public/")
	s.HandleFunc("GET", "/login", func(a *rego.App) rego.Result {
		return a.RenderTemplate("/public/login.html")
	})
	s.HandleFunc("POST", "/login", func(a *rego.App) rego.Result {
		fmt.Println(a.Params["username"])
		fmt.Println(a.Params["password"])
		if a.Params["username"] == "test" && a.Params["password"] == "password" {
			http.SetCookie(a.ResponseWriter, &http.Cookie{
				Name:  "X_AUTH",
				Value: Sign("verified"),
				Path:  "/",
			})
			return a.Redirect("/public/index.html")
		}
		return a.RenderTemplate("/public/login.html")

	})

	s.Use(AuthHandler)

	s.Run(":8082")
}

func Index(a *rego.App) rego.Result {
	return "Welcome rego"
}

func PostUser(a *rego.App) rego.Result {
	if u, ok := a.Params["user"]; ok {
		log.Println(u)
	}
	return nil
}

func Users(a *rego.App) rego.Result {
	users := []User{User{1, "John", "john@mail.com"}, User{2, "Bob", "bob@mail.com"}, User{3, "Mark", "mark@mail.com"}}
	return a.RenderJson(users)
}

func AuthHandler(next http.Handler) http.Handler {
	ignore := []string{"/login"}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, s := range ignore {
			if strings.HasPrefix(r.URL.Path, s) {
				next.ServeHTTP(w, r)
				return
			}
		}

		if v, err := r.Cookie("X_AUTH"); err == http.ErrNoCookie {
			// not authenticated
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else if err != nil {
			// some other error
			panic(err.Error())
		} else if Verify("verified", v.Value) {
			// success - call the next handler
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
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
