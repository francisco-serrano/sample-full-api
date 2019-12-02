package utils

import (
	"fmt"
	"net/http"
)

type ForecastError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

const (
	ErrInvalidData         = "err_invalid_data"
	ErrNotFound            = "err_not_found"
	ErrInternalServerError = "err_internal_server_error"
)

func (e *ForecastError) Error() string {
	return fmt.Sprintf("[%d][%s]:%s", GetStatusErrorCode(e), e.Code, e.Description)
}

func ErrorInvalid(description string) *ForecastError {
	return &ForecastError{
		Code:        ErrInvalidData,
		Description: description,
	}
}

func ErrorNotFound(description string) *ForecastError {
	return &ForecastError{
		Code:        ErrNotFound,
		Description: description,
	}
}

func ErrorInternal(description string) *ForecastError {
	return &ForecastError{
		Code:        ErrInternalServerError,
		Description: description,
	}
}

func GetStatusErrorCode(error *ForecastError) int {
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
