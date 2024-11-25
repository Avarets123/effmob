package middlewares

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CORS(next httprouter.Handle) httprouter.Handle {
	return (func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		origin := r.Header.Get("Origin")

		if origin == "" {
			next(w, r, p)
			return
		}

		headers := w.Header()

		headers.Set("Access-Control-Allow-Origin", origin)
		headers.Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			headers.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD,PATCH,OPTIONS")
			headers.Set("Access-Control-Allow-Headers", "authorization,content-type,content-length")
			headers.Set("Access-Control-Max-Age", "86400")
		}

		next(w, r, p)

	})
}
