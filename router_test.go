package rego

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRouter(t *testing.T) {

	Convey("Get", t, func() {
		r := &router{dispatchers: make(map[string]*dispatcher)}

		ok := false
		r.HandleFunc("GET", "/users", func(c *Context) Result { ok = true; return nil })

		req, _ := http.NewRequest("GET", "/users", nil)
		r.ServeHTTP(&mockResponseWriter{}, req)

		So(ok, ShouldBeTrue)

	})

	Convey("Post", t, func() {
		r := &router{dispatchers: make(map[string]*dispatcher)}

		ok := false
		r.HandleFunc("POST", "/users", func(c *Context) Result { ok = true; return nil })

		req, _ := http.NewRequest("POST", "/users", nil)
		r.ServeHTTP(&mockResponseWriter{}, req)

		So(ok, ShouldBeTrue)
	})

}
