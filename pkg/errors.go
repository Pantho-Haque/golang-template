package pkg

import (
	"fmt"
	"maps"
	"net/http"
)

// StatusCoder defines status code for http code
type StatusCoder interface {
	StatusCode() int
}

// HTTPFormatter formates http error
type HTTPFormatter interface {
	HTTPFormat() map[string]any
}

// Error defines error
type Error struct {
	LogMsg    string
	ClientMsg any
	Code      int
}

func (e *Error) Error() string {
	return fmt.Sprintf("log msg: %s, status code: %d", e.LogMsg, e.Code)
}

func (e *Error) HTTPFormat() map[string]any {
	m := map[string]any{
		"status": false,
	}
	switch msg := e.ClientMsg.(type) {
	case map[string]any:
		maps.Copy(m, msg)
	default:
		m["error"] = msg
	}
	return m
}

func (e *Error) StatusCode() int {
	return e.Code
}

// NewBadReq returns error for bad request
func NewBadReq(logMsg string, clientMsg any) error {
	return NewClientError(logMsg, clientMsg, http.StatusBadRequest)
}

// NewForbidden returns error for forbidden
func NewForbidden(logMsg, clientMsg string) error {
	return NewClientError(logMsg, clientMsg, http.StatusForbidden)
}

func NewConflict(logMsg, clientMsg string) error {
	return NewClientError(logMsg, clientMsg, http.StatusConflict)
}

func NewNotAcceptable(logMsg, clientMsg string) error {
	return NewClientError(logMsg, clientMsg, http.StatusNotAcceptable)
}

// NewNotFound returns error for not found
func NewNotFound(logMsg, clientMsg string) error {
	return NewClientError(logMsg, clientMsg, http.StatusNotFound)
}

// NewInternalServer returns error for internal server
func NewInternalServer(logMsg, clientMsg string) error {
	if clientMsg == "" {
		clientMsg = "internal server error"
	}
	return NewClientError(logMsg, clientMsg, http.StatusInternalServerError)
}

func NewUnprocessable(logMsg, clientMsg string) error {
	return NewClientError(logMsg, clientMsg, http.StatusUnprocessableEntity)
}

func NewClientError(logMsg string, clientMsg any, code int) error {
	return &Error{
		LogMsg:    logMsg,
		ClientMsg: clientMsg,
		Code:      code,
	}
}
