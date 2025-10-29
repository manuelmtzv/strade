package transport

import (
	"net/http"
	"strade/internal/constants"
	"strade/internal/utils"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (t *Transporter) InternalServerError(w http.ResponseWriter, r *http.Request, err error, messageID string) {
	t.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	msg := t.translator.GetMessageOrDefault(r.Context(), messageID, "internal server error", nil)

	t.WriteJSONError(w, http.StatusInternalServerError, msg)
}

func (t *Transporter) InternalServerErrorBasic(w http.ResponseWriter, r *http.Request, err error) {
	t.InternalServerError(w, r, err, constants.ErrorInternalServerError)
}

func (t *Transporter) ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	t.logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path)

	t.WriteJSONError(w, http.StatusForbidden, t.translator.GetMessageOrDefault(r.Context(), constants.ErrorForbidden, "this action is forbidden", nil))
}

func (t *Transporter) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	t.logger.Warnw("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	t.WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (t *Transporter) ConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	msg, localizeErr := t.translator.GetLocaleMessage(r.Context(), constants.ErrorConflict, nil)
	if err != nil && localizeErr == nil {
		msg = err.Error()
	}

	t.logger.Errorw("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	t.WriteJSONError(w, http.StatusConflict, msg)
}

func (t *Transporter) NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	msg := "not found"
	if err != nil {
		msg = err.Error()
	}

	t.logger.Warnw("not found", "method", r.Method, "path", r.URL.Path, "error", msg)

	t.WriteJSONError(w, http.StatusNotFound, msg)
}

func (t *Transporter) UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	t.logger.Warnw("unauthorized", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	t.WriteJSONError(w, http.StatusUnauthorized, err.Error())
}

func (t *Transporter) UnauthorizedGenericErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	t.logger.Warnw("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	t.WriteJSONError(w, http.StatusUnauthorized, t.translator.GetMessageOrDefault(r.Context(), constants.ErrorUnauthorized, "unauthorized", nil))
}

func (t *Transporter) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	t.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	msg := "rate limit exceeded, retry after: " + retryAfter
	if localizer, err := utils.GetLocalizerFromContext(r.Context()); err == nil {
		msg, _ = localizer.Localize(&i18n.LocalizeConfig{
			MessageID:    constants.ErrorRateLimitExceeded,
			TemplateData: map[string]string{"RetryAfter": retryAfter},
		})
	}

	t.WriteJSONError(w, http.StatusTooManyRequests, msg)
}
