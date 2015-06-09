package rego

import (
	"net/http"
	"sync"
)

var ctx *context

func init() {
	ctx = &context{data: make(map[*http.Request]map[string]interface{})}
}

type context struct {
	mutex sync.RWMutex
	data  map[*http.Request]map[string]interface{}
}

func (c *context) Set(r *http.Request, k string, v interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.data[r] == nil {
		c.data[r] = make(map[string]interface{})
	}
	c.data[r][k] = v
}

func (c *context) Get(r *http.Request, k string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if ctx, ok := c.data[r]; ok {
		v, ok := ctx[k]
		return v, ok
	}
	return nil, false
}

func (c *context) GetAll(r *http.Request) (map[string]interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ctx, ok := c.data[r]
	result := make(map[string]interface{}, len(ctx))
	for k, v := range ctx {
		result[k] = v
	}

	return result, ok
}

func (c *context) Delete(r *http.Request, k string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.data[r] != nil {
		delete(c.data[r], k)
	}
}

func (c *context) Clear(r *http.Request) {
	c.mutex.Lock()
	c.mutex.Unlock()

	delete(c.data, r)
}

func (c *context) ClearHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer c.Clear(r)
		h.ServeHTTP(w, r)
	})
}
