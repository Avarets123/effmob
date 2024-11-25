package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func LogMid(logger *slog.Logger) Middleware {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			start := time.Now()
			newWriter := &WriterWithStatus{
				ResponseWriter: w,
			}
			next(newWriter, r, p)
			logStr := fmt.Sprintf("%s %s; duration: %s; statusCode: %d; remoteAdd: %s;", r.Method, r.URL, time.Since(start), newWriter.StatusCode, r.RemoteAddr)
			logger.Info(logStr)

		}
	}
}
