package errors

import (
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// ErrorRsp
/*
Response code and description error
*/
type ErrorRsp struct {
	ErrorCode string `json:"ErrorCode" xml:"ErrorCode"`
	ErrorDesc string `json:"ErrorDesc" xml:"ErrorDesc"`
}

func (e *ErrorRsp) Error() string {
	return e.ErrorDesc
}

func NewCustomErrorRsp(code int, msg string) *echo.HTTPError {
	if len(msg) == 0 {
		msg = http.StatusText(code)
	}

	err := echo.NewHTTPError(code, ErrorRsp{
		ErrorCode: strconv.Itoa(code),
		ErrorDesc: msg,
	})

	return err.SetInternal(errors.WithStack(errors.New(msg)))
}

// NewInternalErrorRsp status500
func NewInternalErrorRsp(msg string) *echo.HTTPError {
	return NewCustomErrorRsp(http.StatusInternalServerError, msg)
}

// NewBadRequestErrorRsp status400
func NewBadRequestErrorRsp(msg string) *echo.HTTPError {
	return NewCustomErrorRsp(http.StatusBadRequest, msg)
}

// NewNotFoundErrorRsp status404
func NewNotFoundErrorRsp(msg string) *echo.HTTPError {
	return NewCustomErrorRsp(http.StatusNotFound, msg)
}

// NewNotAcceptableErrorRsp status406
func NewNotAcceptableErrorRsp(msg string) *echo.HTTPError {
	return NewCustomErrorRsp(http.StatusNotAcceptable, msg)
}

// NewUnsupportedMediaTypeErrRsp status415
func NewUnsupportedMediaTypeErrRsp(msg string) *echo.HTTPError {
	return NewCustomErrorRsp(http.StatusUnsupportedMediaType, msg)
}

// NewRangeNotSatisfiableErrRsp status416
func NewRangeNotSatisfiableErrRsp(msg string) *echo.HTTPError {
	return NewCustomErrorRsp(http.StatusRequestedRangeNotSatisfiable, msg)
}

// NewClientRestErrRsp custom status
func NewClientRestErrRsp(code int, body io.Reader) *echo.HTTPError {
	var desc string
	all, err := io.ReadAll(body)
	desc = string(all)
	if err != nil {
		desc = err.Error()
	}
	return NewCustomErrorRsp(code, desc)
}

type UnmarshalError struct {
	Err ErrorRsp
}

func (e *UnmarshalError) Error() string {
	return e.Err.ErrorDesc
}

func (e *UnmarshalError) Unwrap() error {
	return e
}

func NewUnmarshalError(err ErrorRsp) error {
	return errors.WithStack(&UnmarshalError{Err: err})
}

type ValidationError struct {
	Err ErrorRsp
}

func (e *ValidationError) Error() string {
	return e.Err.ErrorDesc
}

func (e *ValidationError) Unwrap() error {
	return e
}

func NewValidationError(err ErrorRsp) error {
	return errors.WithStack(&ValidationError{Err: err})
}

type InternalError struct {
	Err ErrorRsp
}

func (e *InternalError) Error() string {
	return e.Err.ErrorDesc
}

func NewInternalError(err ErrorRsp) error {
	return errors.WithStack(&InternalError{Err: err})
}

type NotFoundError struct {
	Err ErrorRsp
}

func (e *NotFoundError) Error() string {
	return e.Err.ErrorDesc
}

func NewNotFoundError(err ErrorRsp) error {
	return errors.WithStack(&NotFoundError{Err: err})
}

type ParametersMissingError struct {
	Err ErrorRsp
}

func (e *ParametersMissingError) Error() string {
	return e.Err.ErrorDesc
}

func NewParametersMissingError(err ErrorRsp) error {
	return errors.WithStack(&ParametersMissingError{Err: err})
}

type BadRequestError struct {
	Err ErrorRsp
}

func (e *BadRequestError) Error() string {
	return e.Err.ErrorDesc
}

func NewBadRequestError(err ErrorRsp) error {
	return errors.WithStack(&BadRequestError{Err: err})
}

type ResponseClient struct {
	Err  ErrorRsp
	Code int
}

func (e *ResponseClient) Error() string {
	return e.Err.ErrorDesc
}

func NewResponseClient(code int, body io.Reader) error {
	var desc string
	all, err := io.ReadAll(body)
	desc = string(all)
	if err != nil {
		desc = err.Error()
	}
	return &ResponseClient{
		Err: ErrorRsp{
			ErrorCode: http.StatusText(code),
			ErrorDesc: desc,
		},
		Code: code,
	}
}
