package service

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type (
	// ErrCode type
	ErrCode uint16
	// Error struct
	Error struct {
		Code           ErrCode
		Err            error
		httpStatusCode int
	}
)

// Error codes
const (
	UnknownErrCode              ErrCode = iota
	NotFoundErrCode                     // resource not found on database
	PermissionNotGrantedErrCode         // user din't granted all the permissions we want
	NoMoreCreditErrCode
	FacebookAPIErrCode
	JWTTokenErrCode
	ValidationErrCode
	AlreadyExistsErrCode
)

//NewErr instances new error
func NewErr(c ErrCode, err error) *Error {
	return &Error{Code: c, Err: err}
}

func (e Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("service error: %d", e.Code)
}

func (e Error) HTTPStatusCode() int {
	if e.httpStatusCode == 0 {
		return http.StatusBadRequest
	}
	return e.httpStatusCode
}

func (e *Error) SetHTTPStatusCode(c int) *Error {
	e.httpStatusCode = c
	return e
}

func isNotFoundErr(err error) bool {
	if e, ok := errors.Cause(err).(*Error); ok && e.Code == NotFoundErrCode {
		return true
	}
	return false
}
