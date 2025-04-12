package errors

import (
	"errors"
	"net/http"
)

type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(status int, code, message string) APIError {
	return APIError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func NewAPIErrorWithDetails(status int, code, message, details string) APIError {
	return APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

var (
	ErrInvalidCredentials = APIError{
		Status:  http.StatusUnauthorized,
		Code:    "INVALID_CREDENTIALS",
		Message: "Invalid email or password",
	}

	ErrDuplicateResource = APIError{
		Status:  http.StatusConflict,
		Code:    "DUPLICATE_RESOURCE",
		Message: "Resource already exists",
	}

	ErrDuplicateEmail = APIError{
		Status:  http.StatusConflict,
		Code:    "DUPLICATE_RESOURCE",
		Message: "Email already exists",
	}

	ErrInvalidID = APIError{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_ID",
		Message: "Invalid ID format",
	}

	ErrNotFound = APIError{
		Status:  http.StatusNotFound,
		Code:    "NOT_FOUND",
		Message: "Resource not found",
	}

	ErrInternalServerError = APIError{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "An unexpected error occurred",
	}
)

func NewInvalidCredentialsError() error {
	return ErrInvalidCredentials
}

func NewDuplicateResourceError() error {
	return ErrDuplicateResource
}

func NewInvalidIDError() error {
	return ErrInvalidID
}

func NewNotFoundError() error {
	return ErrNotFound
}

func NewInternalServerError() error {
	return ErrInternalServerError
}

func IsAPIError(err error) bool {
	var apiErr APIError
	return errors.As(err, &apiErr)
}
