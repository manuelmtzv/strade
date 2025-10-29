package transport

import (
	"net/http"
	"strade/internal/constants"
	"strade/internal/utils"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ErrorOption func(*errorConfig)

type errorConfig struct {
	technicalError error
	messageID      string
	message        string
	templateData   map[string]any
	errorCode      string
}

func WithTechnicalError(err error) ErrorOption {
	return func(c *errorConfig) {
		c.technicalError = err
	}
}

func WithMessageID(messageID string) ErrorOption {
	return func(c *errorConfig) {
		c.messageID = messageID
	}
}

func WithMessage(message string) ErrorOption {
	return func(c *errorConfig) {
		c.message = message
	}
}

func WithTemplateData(data map[string]any) ErrorOption {
	return func(c *errorConfig) {
		c.templateData = data
	}
}

func WithErrorCode(errorCode string) ErrorOption {
	return func(c *errorConfig) {
		c.errorCode = errorCode
	}
}

func (t *Transporter) sendError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	logLevel string,
	logMessage string,
	defaultMessage string,
	opts ...ErrorOption,
) {
	cfg := &errorConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	logFields := []any{"method", r.Method, "path", r.URL.Path}
	if cfg.technicalError != nil {
		logFields = append(logFields, "error", cfg.technicalError.Error())
	}

	switch logLevel {
	case "error":
		t.logger.Errorw(logMessage, logFields...)
	case "warn":
		t.logger.Warnw(logMessage, logFields...)
	default:
		t.logger.Infow(logMessage, logFields...)
	}

	userMessage := defaultMessage
	if cfg.message != "" {
		userMessage = cfg.message
	} else if cfg.messageID != "" {
		userMessage = t.translator.GetMessageOrDefault(r.Context(), cfg.messageID, defaultMessage, cfg.templateData)
	}

	errorCode := cfg.errorCode
	if errorCode == "" {
		errorCode = "UNKNOWN_ERROR"
	}

	t.WriteJSONError(w, status, userMessage, errorCode)
}

func (t *Transporter) InternalServerError(w http.ResponseWriter, r *http.Request, opts ...ErrorOption) {
	t.sendError(w, r, http.StatusInternalServerError, "error", "internal server error", "internal server error", opts...)
}

func (t *Transporter) BadRequest(w http.ResponseWriter, r *http.Request, opts ...ErrorOption) {
	t.sendError(w, r, http.StatusBadRequest, "warn", "bad request", "bad request", opts...)
}

func (t *Transporter) NotFound(w http.ResponseWriter, r *http.Request, opts ...ErrorOption) {
	t.sendError(w, r, http.StatusNotFound, "warn", "not found", "not found", opts...)
}

func (t *Transporter) Conflict(w http.ResponseWriter, r *http.Request, opts ...ErrorOption) {
	t.sendError(w, r, http.StatusConflict, "error", "conflict", "conflict", opts...)
}

func (t *Transporter) Unauthorized(w http.ResponseWriter, r *http.Request, opts ...ErrorOption) {
	t.sendError(w, r, http.StatusUnauthorized, "warn", "unauthorized", "unauthorized", opts...)
}

func (t *Transporter) Forbidden(w http.ResponseWriter, r *http.Request, opts ...ErrorOption) {
	t.sendError(w, r, http.StatusForbidden, "warn", "forbidden", "this action is forbidden", opts...)
}

func (t *Transporter) InternalServerErrorBasic(w http.ResponseWriter, r *http.Request, err error) {
	t.InternalServerError(w, r, WithTechnicalError(err), WithMessageID(constants.ErrorInternalServerError))
}

func (t *Transporter) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	t.BadRequest(w, r, WithTechnicalError(err), WithMessage(err.Error()))
}

func (t *Transporter) ConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	msg, localizeErr := t.translator.GetLocaleMessage(r.Context(), constants.ErrorConflict, nil)
	if err != nil && localizeErr == nil {
		msg = err.Error()
	}
	t.Conflict(w, r, WithTechnicalError(err), WithMessage(msg))
}

func (t *Transporter) NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	msg := "not found"
	if err != nil {
		msg = err.Error()
	}
	t.NotFound(w, r, WithTechnicalError(err), WithMessage(msg))
}

func (t *Transporter) UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	t.Unauthorized(w, r, WithTechnicalError(err), WithMessage(err.Error()))
}

func (t *Transporter) ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	t.Forbidden(w, r, WithMessageID(constants.ErrorForbidden))
}

func (t *Transporter) UnauthorizedGenericErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	t.Unauthorized(w, r, WithTechnicalError(err), WithMessageID(constants.ErrorUnauthorized))
}

func (t *Transporter) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	w.Header().Set("Retry-After", retryAfter)

	msg := "rate limit exceeded, retry after: " + retryAfter
	if localizer, err := utils.GetLocalizerFromContext(r.Context()); err == nil {
		msg, _ = localizer.Localize(&i18n.LocalizeConfig{
			MessageID:    constants.ErrorRateLimitExceeded,
			TemplateData: map[string]string{"RetryAfter": retryAfter},
		})
	}

	t.WriteJSONError(w, http.StatusTooManyRequests, msg, constants.ErrCodeRateLimitExceeded)
}
