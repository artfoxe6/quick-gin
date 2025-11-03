package apperr

import (
	"errors"
	"fmt"
	"net/http"
)

// Error represents a structured application error with an associated HTTP status code.
type Error struct {
	Code    int
	Message string
	Err     error
}

// New creates a new application error with the given status code and message.
func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

// Wrap creates a new application error with an underlying error for context.
func Wrap(code int, message string, err error) *Error {
	if err == nil {
		return New(code, message)
	}
	return &Error{Code: code, Message: message, Err: err}
}

// Internal returns an application error representing an internal server error.
func Internal(err error) *Error {
	if err == nil {
		err = errors.New(http.StatusText(http.StatusInternalServerError))
	}
	return Wrap(http.StatusInternalServerError, err.Error(), err)
}

// BadRequest creates a bad request application error.
func BadRequest(message string) *Error {
	return New(http.StatusBadRequest, message)
}

// Unauthorized creates an unauthorized application error.
func Unauthorized(message string) *Error {
	return New(http.StatusUnauthorized, message)
}

// Forbidden creates a forbidden application error.
func Forbidden(message string) *Error {
	return New(http.StatusForbidden, message)
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("app error (code=%d)", e.Code)
}

// Unwrap exposes the wrapped error for compatibility with errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
