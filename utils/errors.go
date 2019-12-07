package utils

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

const (
	ErrInvalidData         = "err_invalid_data"
	ErrUnauthorized        = "err_unauthorized"
	ErrNotFound            = "err_not_found"
	ErrInternalServerError = "err_internal_server_error"
)

func (e *Error) Error() string {
	return fmt.Sprintf("[%d][%s]:%s", GetStatusErrorCode(e), e.Code, e.Description)
}

func ErrorInvalid(description string) *Error {
	return &Error{
		Code:        ErrInvalidData,
		Description: description,
	}
}

func ErrorUnauthorized(description string) *Error {
	return &Error{
		Code:        ErrUnauthorized,
		Description: description,
	}
}

func ErrorNotFound(description string) *Error {
	return &Error{
		Code:        ErrNotFound,
		Description: description,
	}
}

func ErrorInternal(description string) *Error {
	return &Error{
		Code:        ErrInternalServerError,
		Description: description,
	}
}

func GetStatusErrorCode(error *Error) int {
	switch error.Code {
	case ErrInvalidData:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInternalServerError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
