package rego

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRouter(t *testing.T) {

	Convey("Get", t, func() {
		r := &router{make(map[string]map[string]HandlerFunc)}

		ok := false
		r.HandleFunc("GET", "/users", func(c *Context) { ok = true })

		req, _ := http.NewRequest("GET", "/users", nil)
		c := &Context{ResponseWriter: &mockResponseWriter{}, Request: req}

		r.handler()(c)

		So(ok, ShouldBeTrue)

	})

	Convey("Post", t, func() {
		r := &router{make(map[string]map[string]HandlerFunc)}

		ok := false
		r.HandleFunc("POST", "/users", func(c *Context) { ok = true })

		req, _ := http.NewRequest("POST", "/users", nil)
		c := &Context{ResponseWriter: &mockResponseWriter{}, Request: req}

		r.handler()(c)

		So(ok, ShouldBeTrue)
	})

	Convey("Lookup", t, func() {
		r := &router{make(map[string]map[string]HandlerFunc)}

		ok := false
		r.HandleFunc("GET", "/users/:id/addresses", func(c *Context) {
			if c.Params["id"] == "1" {
				ok = true
			}
		})

		Convey("found", func() {
			req, _ := http.NewRequest("GET", "/users/1/addresses", nil)
			c := &Context{
				Params:         make(map[string]interface{}),
				ResponseWriter: &mockResponseWriter{},
				Request:        req,
			}
			r.handler()(c)

			So(ok, ShouldBeTrue)
		})

		Convey("not found", func() {
			req, _ := http.NewRequest("GET", "/users/2/addresses", nil)
			c := &Context{
				Params:         make(map[string]interface{}),
				ResponseWriter: &mockResponseWriter{},
				Request:        req,
			}
			r.handler()(c)

			So(ok, ShouldBeFalse)
		})

	})
}
