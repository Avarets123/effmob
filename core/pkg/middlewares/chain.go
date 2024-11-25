package middlewares

import (
	"github.com/julienschmidt/httprouter"
)

func Chain(middlewares ...Middleware) Middleware {
	return func(h httprouter.Handle) httprouter.Handle {
		for _, mdw := range middlewares {
			h = mdw(h)
		}
		return h
	}

}
