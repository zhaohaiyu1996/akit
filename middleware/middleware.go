package middleware

import "context"

// middleware function
type MiddleWareFunc func(context.Context, interface{}) (interface{}, error)

// use by construct middleware function
type MiddleWare func(MiddleWareFunc) MiddleWareFunc

// Chain is a middleware Chain by all middleware and handle
func Chain(others ...MiddleWare) MiddleWare {
	return func(next MiddleWareFunc) MiddleWareFunc {
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return next
	}
}
