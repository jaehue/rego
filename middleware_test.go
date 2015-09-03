package rego

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRecover(t *testing.T) {
	Convey("Get", t, func() {
		r := &router{
			dispatchers:       make(map[string]*dispatcher),
			staticFileHandler: make(map[string]http.HandlerFunc),
		}
		r.HandleFunc("GET", "/", func(a *App) Result { panic("panic!") })

		handler := recoverHandler(r)
		req, _ := http.NewRequest("GET", "/", nil)

		So(func() { handler.ServeHTTP(&mockResponseWriter{}, req) }, ShouldNotPanic)
	})
}
