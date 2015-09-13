package rego

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRecover(t *testing.T) {
	Convey("Get", t, func() {
		r := &router{make(map[string]map[string]HandlerFunc)}
		r.HandleFunc("GET", "/", func(c *Context) { panic("panic!") })

		handler := recoverHandler(r.handler())
		req, _ := http.NewRequest("GET", "/", nil)
		c := &Context{ResponseWriter: &mockResponseWriter{}, Request: req}

		So(func() { handler(c) }, ShouldNotPanic)
	})
}
