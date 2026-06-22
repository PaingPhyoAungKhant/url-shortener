package apperr_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/apperr"
)

func TestConstructors_SetCodeMessageAndCause(t *testing.T) {
	cause := errors.New("db connection refused")

	tests := []struct {
		name      string
		err       *apperr.AppErr
		wantCode  apperr.AppErrCode
		wantMsg   string
		wantCause error
	}{
		{
			name:      "Not found with cause",
			err:       apperr.NotFound("user not found", cause),
			wantCode:  apperr.AppErrNotFound,
			wantMsg:   "user not found",
			wantCause: cause,
		},
		{
			name:      "Not found without cause",
			err:       apperr.NotFound("user not found", nil),
			wantCode:  apperr.AppErrNotFound,
			wantMsg:   "user not found",
			wantCause: nil,
		},
		{
			name:      "Unauthorized",
			err:       apperr.UnAuthorized("invalid token", cause),
			wantCode:  apperr.AppErrUnAuthorized,
			wantMsg:   "invalid token",
			wantCause: cause,
		},
		{
			name:      "Forbidden",
			err:       apperr.Forbidden("access denied", cause),
			wantCode:  apperr.AppErrForbidden,
			wantMsg:   "access denied",
			wantCause: cause,
		},
		{
			name:      "Bad Request",
			err:       apperr.BadRequest("invalid email", cause),
			wantCode:  apperr.AppErrBadRequest,
			wantMsg:   "invalid email",
			wantCause: cause,
		},
		{
			name:      "Internal",
			err:       apperr.Internal("something went wrong", cause),
			wantCode:  apperr.AppErrInternal,
			wantMsg:   "something went wrong",
			wantCause: cause,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantCode, tc.err.Code)
			assert.Equal(t, tc.wantMsg, tc.err.Message)
			assert.Equal(t, tc.wantCause, tc.err.Unwrap())
		})
	}
}

func TestError_FormatAsCodeColonMessage(t *testing.T) {
	msg := "user not found"
	err := apperr.NotFound(msg, nil)
	want := fmt.Sprintf("%s: %s", apperr.AppErrNotFound, msg)

	assert.Equal(t, want, err.Error())
}

func TestLogError_WithCauseIncludeCause(t *testing.T) {
	cause := errors.New("db connection refused")
	msg := "user not found"
	err := apperr.NotFound(msg, cause)
	want := fmt.Sprintf("%s: %s: %s", apperr.AppErrNotFound, msg, cause)
	assert.Equal(t, want, err.LogError())
}

func TestLogError_WithoutCauseEqualError(t *testing.T) {
	err := apperr.Internal("something went wrong", nil)
	assert.Equal(t, err.Error(), err.LogError())
}

func TestIs_SameCodeMatches(t *testing.T) {
	err := apperr.NotFound("user not found", nil)
	target := apperr.ErrNotFound
	matches := errors.Is(err, target)
	want := true
	assert.Equal(t, want, matches)
}

func TestIs_DifferentCode_DoesNotMatch(t *testing.T) {
	err := apperr.BadRequest("invalid email", nil)
	target := apperr.ErrNotFound
	matches := errors.Is(err, target)
	want := false
	assert.Equal(t, want, matches)
}

func TestIs_EmptyCode_MatchesAnyAppErr(t *testing.T) {
	err := apperr.NotFound("user not found", nil)
	assert.True(t, errors.Is(err, &apperr.AppErr{}))
}

func TestIs_NonAppErrTarget_ReturnsFalse(t *testing.T) {
	err := errors.New("db connection refused")
	target := &apperr.AppErr{}
	matches := errors.Is(err, target)
	want := false
	assert.Equal(t, want, matches)
}

func TestIs_WrappedAppErr_FoundThroughChain(t *testing.T) {
	cause := apperr.NotFound("slug not found", nil)
	err := apperr.Internal("resolve failed", cause)

	assert.True(t, errors.Is(err, apperr.ErrInternal))
	assert.True(t, errors.Is(err, apperr.ErrNotFound))
}

func TestWithAppErr_AppErrFromCtx_RoundTrip(t *testing.T) {
	ctx := apperr.WithAppErr(context.Background())
	apperr.AppErrSetErrFromCtx(ctx, apperr.ErrNotFound)
	err := apperr.AppErrFromCtx(ctx)
	matches := apperr.ErrNotFound.Is(err)
	want := true
	assert.Equal(t, want, matches)
}

func TestAppErrFromCtx_ReturnNilWhenNoSlot(t *testing.T) {
	err := apperr.AppErrFromCtx(context.Background())
	assert.Nil(t, err)
}

func TestAppErrFromCtx_ReturnNilWhenSlotUnset(t *testing.T) {
	ctx := apperr.WithAppErr(context.Background())
	err := apperr.AppErrFromCtx(ctx)
	assert.Nil(t, err)
}
