package utils

import (
	"fmt"
	"log/slog"
	"net/http"
)

func LogUnhandlingRequest(r *http.Request, logger *slog.Logger) {
	logger.Info(
		fmt.Sprintf("Method: %v, path: %v, remoteAdd: %v, authHeader: %v", r.Method, r.URL, r.RemoteAddr, r.Header.Get("Authorization")),
	)

}
