package rego

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRouter(t *testing.T) {

	s := New()

	Convey("Get", t, func() {

		ok := false

		s.Get("/users", func(c *Context) Result { ok = true; return nil })

		r, _ := http.NewRequest("GET", "/users", nil)
		s.mux.ServeHTTP(&mockResponseWriter{}, r)

		So(ok, ShouldBeTrue)
	})

	Convey("Post", t, func() {
		ok := false

		s.Post("/users", func(c *Context) Result { ok = true; return nil })

		r, _ := http.NewRequest("GET", "/users", nil)
		s.mux.ServeHTTP(&mockResponseWriter{}, r)

		So(ok, ShouldBeTrue)
	})

}
