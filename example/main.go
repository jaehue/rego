package main

import "github.com/jaehue/rego"

func main() {
	s := rego.New()
	s.Get("/", Index)
	s.Run(":8080")
}

func Index(c *rego.Context) rego.Result {
	return "hello rego"
}
