package apperror

import "errors"

type Code string

const (
	CodeInvalidPath      Code = "invalid_path"
	CodeInvalidConfig    Code = "invalid_config"
	CodeOutsideRoot      Code = "outside_root"
	CodeIsDirectory      Code = "is_directory"
	CodeNotDirectory     Code = "not_directory"
	CodeNotFound         Code = "not_found"
	CodePermissionDenied Code = "permission_denied"
	CodeFilesystem       Code = "filesystem"
	CodeRequestCanceled  Code = "request_canceled"
	CodeRequestTimeout   Code = "request_timeout"
	CodeServer           Code = "server"
)

type OperationalError struct {
	Code    Code
	Message string
	Err     error
}

func New(code Code, message string) *OperationalError {
	return &OperationalError{
		Code:    code,
		Message: message,
	}
}

func Wrap(
	code Code,
	message string,
	err error,
) *OperationalError {
	return &OperationalError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func WrapOperation(
	code Code,
	message string,
	operation string,
	err error,
) *OperationalError {
	return Wrap(code, message, operationError{
		operation: operation,
		err:       err,
	})
}

func (e *OperationalError) Error() string {
	if e.Err == nil {
		return e.Message
	}

	return e.Message + ": " + e.Err.Error()
}

func (e *OperationalError) Unwrap() error {
	return e.Err
}

func (e *OperationalError) Is(target error) bool {
	targetError, ok := target.(*OperationalError)
	if !ok {
		return false
	}

	return e.Code == targetError.Code
}

func AsOperational(err error) (*OperationalError, bool) {
	var opErr *OperationalError
	if errors.As(err, &opErr) {
		return opErr, true
	}

	return nil, false
}

func HasCode(err error, code Code) bool {
	opErr, ok := AsOperational(err)
	return ok && opErr.Code == code
}

type operationError struct {
	operation string
	err       error
}

func (e operationError) Error() string {
	return e.operation + ": " + e.err.Error()
}

func (e operationError) Unwrap() error {
	return e.err
}
