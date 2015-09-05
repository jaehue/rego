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
		r.HandleFunc("GET", "/users", func(a *App) { ok = true })

		req, _ := http.NewRequest("GET", "/users", nil)
		a := &App{Params: make(map[string]interface{}), ResponseWriter: &mockResponseWriter{}, Request: req}
		r.handle()(a)

		So(ok, ShouldBeTrue)

	})

	Convey("Post", t, func() {
		r := &router{dispatchers: make(map[string]*dispatcher)}

		ok := false
		r.HandleFunc("POST", "/users", func(a *App) { ok = true })

		req, _ := http.NewRequest("POST", "/users", nil)
		a := &App{Params: make(map[string]interface{}), ResponseWriter: &mockResponseWriter{}, Request: req}
		r.handle()(a)

		So(ok, ShouldBeTrue)
	})

	Convey("Lookup", t, func() {
		r := &router{dispatchers: make(map[string]*dispatcher)}

		ok := false
		r.HandleFunc("GET", "/users/:id/addresses", func(a *App) {
			if a.Params["id"] == "1" {
				ok = true
			}
		})

		Convey("found", func() {
			req, _ := http.NewRequest("GET", "/users/1/addresses", nil)
			a := &App{Params: make(map[string]interface{}), ResponseWriter: &mockResponseWriter{}, Request: req}
			r.handle()(a)

			So(ok, ShouldBeTrue)
		})

		Convey("not found", func() {
			req, _ := http.NewRequest("GET", "/users/2/addresses", nil)
			a := &App{Params: make(map[string]interface{}), ResponseWriter: &mockResponseWriter{}, Request: req}
			r.handle()(a)

			So(ok, ShouldBeFalse)
		})

	})
}
