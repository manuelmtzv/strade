package handle

import (
	"net/http"
	"strade/internal/constants"
)

func (h *Handler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	msg := h.Translator.GetMessageOrDefault(r.Context(), constants.HealthOk, "Ok", nil)

	h.Logger.Info(msg)

	if err := h.Transporter.JSONResponse(w, http.StatusOK, map[string]any{"message": msg}); err != nil {
		h.Transporter.InternalServerErrorBasic(w, r, err)
	}
}
