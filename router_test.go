package rego

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRouter(t *testing.T) {

	Convey("Get", t, func() {
		r := &router{mux: http.NewServeMux(), dispatchers: make(map[string]*dispatcher)}

		ok := false
		r.Get("/users", func(c *Context) Result { ok = true; return nil })
		r.setHandler()

		req, _ := http.NewRequest("GET", "/users", nil)
		r.ServeHTTP(&mockResponseWriter{}, req)

		So(ok, ShouldBeTrue)

	})

	Convey("Post", t, func() {
		r := &router{mux: http.NewServeMux(), dispatchers: make(map[string]*dispatcher)}

		ok := false
		r.Post("/users", func(c *Context) Result { ok = true; return nil })
		r.setHandler()

		req, _ := http.NewRequest("POST", "/users", nil)
		r.ServeHTTP(&mockResponseWriter{}, req)

		So(ok, ShouldBeTrue)
	})

}
