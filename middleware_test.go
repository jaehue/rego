package rego

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRecover(t *testing.T) {
	Convey("Get", t, func() {
		r := &router{dispatchers: make(map[string]*dispatcher)}
		r.HandleFunc("GET", "/", func(a *App) { panic("panic!") })

		handler := recoverHandler(r.handle())
		req, _ := http.NewRequest("GET", "/", nil)
		a := &App{Params: make(map[string]interface{}), ResponseWriter: &mockResponseWriter{}, Request: req}

		So(func() { handler(a) }, ShouldNotPanic)
	})
}
