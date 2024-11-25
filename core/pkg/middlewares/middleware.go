package middlewares

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Middleware func(next httprouter.Handle) httprouter.Handle

type WriterWithStatus struct {
	StatusCode int
	http.ResponseWriter
}

func (w *WriterWithStatus) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}
