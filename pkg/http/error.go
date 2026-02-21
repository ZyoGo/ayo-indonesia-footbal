package common

import (
	"errors"
	"net/http"

	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/logger"
)

// DefaultResponse default payload response
type DefaultResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload,omitempty"`
}

// CreatedResponse default payload response
type CreatedSuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Payload string `json:"payload"`
}

// ErrorResponse error response
type ErrorResponse struct {
	Code     int    `json:"code"`
	Status   string `json:"status,omitempty"`
	Message  string `json:"message"`
	Internal error  `json:"-"`
}

// NewBadRequestResponse default not found error response
func NewBadRequestResponse() DefaultResponse {
	return DefaultResponse{
		Code:    400,
		Status:  BadRequestStatus,
		Message: "Invalid Request Body",
	}
}

func NewDuplicateResponse() DefaultResponse {
	return DefaultResponse{
		Code:    400,
		Status:  BadRequestStatus,
		Message: "Duplicate",
	}
}

// NewUnauthorizedResponse default unauthorized response
func NewUnauthorizedResponse(msg string) DefaultResponse {
	return DefaultResponse{
		Code:    401,
		Status:  UnauthorizedStatus,
		Message: msg,
	}
}

// NewForbiddenResponse default forbidden response
func NewForbiddenResponse() DefaultResponse {
	return DefaultResponse{
		Code:    403,
		Status:  ForbiddenStatus,
		Message: "Forbidden",
	}
}

// NewInternalServerErrorResponse default internal server error response
func NewInternalServerErrorResponse() DefaultResponse {
	return DefaultResponse{
		Code:    500,
		Status:  InternalErrStatus,
		Message: "Internal server error",
	}
}

// NewValidationErrorResponse default validation error response
func NewValidationErrorResponse(message string) DefaultResponse {
	return DefaultResponse{
		Code:    400,
		Status:  ValidationErrStatus,
		Message: message,
	}
}

// use this one controller or handle layer
func RenderErrorResponse(err error) (resp ErrorResponse) {
	resp = ErrorResponse{Status: InternalErrStatus, Code: http.StatusInternalServerError, Message: "Internal server error", Internal: err}
	var ierr *derrors.Error
	if !errors.As(err, &ierr) {
		logger.Get().Error("error response", "error", err.Error())
		return resp
	} else {
		logger.Get().Error("error response", "error", ierr.Error())
		switch ierr.Code() {
		case derrors.ErrorCodeBadRequest:
			return ErrorResponse{Status: BadRequestStatus, Code: http.StatusBadRequest, Message: ierr.Message(), Internal: ierr}
		case derrors.ErrorCodeUnauthorized:
			return ErrorResponse{Status: UnauthorizedStatus, Code: http.StatusUnauthorized, Message: ierr.Message(), Internal: ierr}
		case derrors.ErrorCodeForbidden:
			return ErrorResponse{Status: ForbiddenStatus, Code: http.StatusForbidden, Message: ierr.Message(), Internal: ierr}
		case derrors.ErrorCodeNotFound:
			return ErrorResponse{Status: NotFoundStatus, Code: http.StatusNotFound, Message: ierr.Message(), Internal: ierr}
		case derrors.ErrorCodeDuplicate:
			return ErrorResponse{Status: DuplicateStatus, Code: http.StatusBadRequest, Message: ierr.Message(), Internal: ierr}
		case derrors.ErrorCodeAlreadyRegistered:
			return ErrorResponse{Status: BadRequestStatus, Code: http.StatusBadRequest, Message: "Already Registered", Internal: ierr}
		case derrors.ErrorCodeInvalidArgument:
			return ErrorResponse{Status: BadRequestStatus, Code: http.StatusBadRequest, Message: "Bad Request", Internal: ierr}
		case derrors.ErrorCodeCustomBadRequest:
			return ErrorResponse{Status: BadRequestStatus, Code: http.StatusBadRequest, Message: ierr.Message(), Internal: ierr}
		}
	}

	return resp
}
