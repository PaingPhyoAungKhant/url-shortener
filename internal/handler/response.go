package handler

import (
	"encoding/json"
	"net/http"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/apperr"
)

// ErrorResponse define the structure for an error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail define the detail of an error
type ErrorDetail struct {
	Code    int               `json:"code"`
	Status  apperr.AppErrCode `json:"status"`
	Message string            `json:"message"`
}

// statusMap map application error code(AppErrCode) to http status code
var statusMap = map[apperr.AppErrCode]int{
	apperr.AppErrNotFound:     http.StatusNotFound,
	apperr.AppErrUnAuthorized: http.StatusUnauthorized,
	apperr.AppErrBadRequest:   http.StatusBadRequest,
	apperr.AppErrForbidden:    http.StatusForbidden,
	apperr.AppErrInternal:     http.StatusInternalServerError,
}

// WriteJSON return a json response with the given statusCode and json body
func WriteJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}

// WriteErr update the appErr value in context for logger and
// return an structured error response
func WriteErr(w http.ResponseWriter, r *http.Request, err *apperr.AppErr) {
	ctx := r.Context()
	if appErr, ok := ctx.Value(apperr.AppErrKey).(**apperr.AppErr); ok {
		*appErr = err
	}
	statusCode, exist := statusMap[err.Code]
	if !exist {
		statusCode = http.StatusInternalServerError
	}
	res := ErrorResponse{
		Error: ErrorDetail{
			Code:    statusCode,
			Status:  err.Code,
			Message: err.Error(),
		},
	}
	WriteJSON(w, statusCode, res)
}
