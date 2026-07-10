package apperr

import "net/http"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *AppError) Error() string { return e.Message }

func NotFound(msg string) *AppError   { return &AppError{http.StatusNotFound, msg} }
func BadRequest(msg string) *AppError { return &AppError{http.StatusBadRequest, msg} }
func Conflict(msg string) *AppError   { return &AppError{http.StatusConflict, msg} }
func Forbidden(msg string) *AppError  { return &AppError{http.StatusForbidden, msg} }
func Internal(msg string) *AppError   { return &AppError{http.StatusInternalServerError, msg} }
