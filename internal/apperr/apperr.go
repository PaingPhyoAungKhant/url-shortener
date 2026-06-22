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

// WithAppErr injects a mutable AppErr slot into ctx so downstream handlers can set it.
// This matches the pattern used by the Error middleware.
func WithAppErr(ctx context.Context) context.Context {
	return context.WithValue(ctx, AppErrKey, new(*AppErr))
}

// AppErrFromCtx returns the AppErr stored in ctx by WithAppErr or the Error middleware,
// or nil if no error has been set.
func AppErrFromCtx(ctx context.Context) *AppErr {
	if appErrPtr, ok := ctx.Value(AppErrKey).(**AppErr); ok && appErrPtr != nil {
		return *appErrPtr
	}
	return nil
}

func AppErrSetErrFromCtx(ctx context.Context, err *AppErr) {
	if appErr, ok := ctx.Value(AppErrKey).(**AppErr); ok {
		*appErr = err
	}
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
