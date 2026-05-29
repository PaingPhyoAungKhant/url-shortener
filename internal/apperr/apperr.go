// Package apperr defines structured application errors with machine-readable
// codes, human-readable messages, and optional root causes for wrapping.
package apperr

import (
	"context"
	"fmt"
)

// AppErrCode is a string enum identifying the category of an application error.
type AppErrCode string

type ctxKey string

const AppErrKey ctxKey = "APP_ERR_KEY"

const (
	AppErrNotFound     AppErrCode = "NOT_FOUND"
	AppErrUnAuthorized AppErrCode = "UNAUTHORIZED"
	AppErrForbidden    AppErrCode = "FORBIDDEN"
	AppErrBadRequest   AppErrCode = "BAD_REQUEST"
	AppErrInternal     AppErrCode = "INTERNAL"
)

// AppErr is a structured error carrying a code, a user-facing message, and an
// optional underlying cause. Use the constructor functions (NotFound, Internal,
// etc.) to build instances rather than constructing the struct directly.
type AppErr struct {
	Code    AppErrCode
	Message string
	cause   error
}

// WithAppErr store the AppErr in given context
func WithAppErr(ctx context.Context, err *AppErr) context.Context {
	return context.WithValue(ctx, AppErrKey, err)
}

// AppErrFromCtx return the AppErr store in given context
// returns nil if AppErr doesn't exist in context
func AppErrFromCtx(ctx context.Context) *AppErr {
	if err, ok := ctx.Value(AppErrKey).(*AppErr); ok {
		return err
	}
	return nil
}

// Error returns a short "CODE: message" string suitable for user-facing output.
func (e *AppErr) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// LogError returns a fuller string that includes the root cause when present,
// intended for structured logging rather than user-facing responses.
func (e *AppErr) LogError() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %s", e.Code, e.Message, e.cause)
	}
	return e.Error()
}

// Unwrap exposes the underlying cause so errors.Is / errors.As can traverse
// the error chain.
func (e *AppErr) Unwrap() error {
	return e.cause
}

// Is reports whether this error matches target. Two AppErr values match when
// they share the same Code, or when target has an empty Code (wildcard match
// against any AppErr regardless of code).
func (e *AppErr) Is(target error) bool {
	t, ok := target.(*AppErr)
	if !ok {
		return false
	}
	if t.Code == "" {
		return true
	}
	return e.Code == t.Code
}

// Sentinel errors for use with errors.Is. They carry only a Code and no
// message, so Is() matches any AppErr with the same Code.
var (
	ErrNotFound     = &AppErr{Code: AppErrNotFound}
	ErrUnAuthorized = &AppErr{Code: AppErrUnAuthorized}
	ErrForbidden    = &AppErr{Code: AppErrForbidden}
	ErrBadRequest   = &AppErr{Code: AppErrBadRequest}
	ErrInternal     = &AppErr{Code: AppErrInternal}
)

// NotFound returns a NOT_FOUND AppErr with the given message and optional cause.
func NotFound(msg string, cause error) *AppErr {
	return &AppErr{
		Code:    AppErrNotFound,
		Message: msg,
		cause:   cause,
	}
}

// UnAuthorized returns an UNAUTHORIZED AppErr with the given message and optional cause.
func UnAuthorized(msg string, cause error) *AppErr {
	return &AppErr{
		Code:    AppErrUnAuthorized,
		Message: msg,
		cause:   cause,
	}
}

// Forbidden returns a FORBIDDEN AppErr with the given message and optional cause.
func Forbidden(msg string, cause error) *AppErr {
	return &AppErr{
		Code:    AppErrForbidden,
		Message: msg,
		cause:   cause,
	}
}

// BadRequest returns a BAD_REQUEST AppErr with the given message and optional cause.
func BadRequest(msg string, cause error) *AppErr {
	return &AppErr{
		Code:    AppErrBadRequest,
		Message: msg,
		cause:   cause,
	}
}

// Internal returns an INTERNAL AppErr with the given message and optional cause.
func Internal(msg string, cause error) *AppErr {
	return &AppErr{
		Code:    AppErrInternal,
		Message: msg,
		cause:   cause,
	}
}
