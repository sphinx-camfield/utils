package stackable

import "net/http"

// Stackable is a function that can be stacked
type Stackable = func(next http.HandlerFunc) http.HandlerFunc

// Stackup stacks multiple stackables into one
func Stackup(stackables ...Stackable) Stackable {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(stackables) - 1; i >= 0; i-- {
			next = stackables[i](next)
		}
		return next
	}
}

// NoopHandlerFunc is doing nothing
func NoopHandlerFunc(_ http.ResponseWriter, _ *http.Request) {}

// HandlerFuncStackup stacks multiple stackables and returns a http.HandlerFunc
func HandlerFuncStackup(stackables ...Stackable) http.HandlerFunc {
	return Stackup(stackables...)(NoopHandlerFunc)
}
