package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strade/internal/translate"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type Error struct {
	Code    string
	Message string
	Details string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) WithDetail(detail string) *Error {
	e.Details = detail
	return e
}

var (
	ErrInternalServer = &Error{
		Code:    "internal_server_error",
		Message: "An internal server error occurred",
	}

	ErrInvalidRequest = &Error{
		Code:    "invalid_request",
		Message: "The request is invalid",
	}

	ErrNotFound = &Error{
		Code:    "not_found",
		Message: "The requested resource was not found",
	}
)

type Transporter struct {
	Validate   *validator.Validate
	logger     *zap.SugaredLogger
	translator translate.Translator
}

func (t *Transporter) SendError(w http.ResponseWriter, r *http.Request, err error) {
	var transportErr *Error
	statusCode := http.StatusInternalServerError

	switch {
	case errors.As(err, &transportErr):
		switch transportErr.Code {
		case ErrInvalidRequest.Code:
			statusCode = http.StatusBadRequest
		case ErrNotFound.Code:
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
	default:
		transportErr = &Error{
			Code:    ErrInternalServer.Code,
			Message: ErrInternalServer.Message,
		}
	}

	if statusCode >= 500 {
		t.logger.Errorw("Internal server error",
			"error", err,
			"path", r.URL.Path,
			"method", r.Method,
		)
	}

	t.SendJSON(w, r, statusCode, ErrorResponse{
		Error: ErrorDetail{
			Code:    transportErr.Code,
			Message: transportErr.Message,
			Details: transportErr.Details,
		},
	})
}

func (t *Transporter) SendJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		t.logger.Errorw("Failed to encode response",
			"error", err,
			"path", r.URL.Path,
			"method", r.Method,
		)
	}
}

func NewTransporter(v *validator.Validate, l *zap.SugaredLogger, t translate.Translator) *Transporter {
	return &Transporter{
		Validate:   v,
		logger:     l,
		translator: t,
	}
}
