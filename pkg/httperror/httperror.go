package httperror

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func New(statusCode int, code, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func BadRequest(message string) *HTTPError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message)
}

func NotFound(message string) *HTTPError {
	return New(http.StatusNotFound, "NOT_FOUND", message)
}

func UnprocessableEntity(message string) *HTTPError {
	return New(http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message)
}

func InternalServerError(message string) *HTTPError {
	return New(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

func WriteError(w http.ResponseWriter, err *HTTPError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(err)
}
