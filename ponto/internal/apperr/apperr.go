package apperr

import "net/http"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Cause }

func NotFound(msg string) *AppError   { return &AppError{Code: http.StatusNotFound, Message: msg} }
func BadRequest(msg string) *AppError { return &AppError{Code: http.StatusBadRequest, Message: msg} }
func Conflict(msg string) *AppError   { return &AppError{Code: http.StatusConflict, Message: msg} }
func Forbidden(msg string) *AppError  { return &AppError{Code: http.StatusForbidden, Message: msg} }

// Internal registra a causa raiz (para log) sem expô-la na resposta HTTP.
func Internal(msg string, cause error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg, Cause: cause}
}
