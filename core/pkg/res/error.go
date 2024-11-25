package res

import "net/http"

type Error interface {
	Error() string
	ResHttp(w http.ResponseWriter)
}

type httpError struct {
	rawError error
	Message  string `json:"message"`
	Status   int    `json:"status"`
	Code     string `json:"code"`
}

func NewError(e error, statusCode int) *httpError {
	return &httpError{
		rawError: e,
		Message:  e.Error(),
		Status:   statusCode,
		Code:     http.StatusText(statusCode),
	}
}

func NewErrorWithMessage(message string, statusCode int) *httpError {
	return &httpError{
		Message: message,
		Status:  statusCode,
		Code:    http.StatusText(statusCode),
	}
}

func (e *httpError) Error() string {
	return e.Message
}

func (e *httpError) ResHttp(w http.ResponseWriter) {
	HttpJsonRes(w, e, e.Status)
}
