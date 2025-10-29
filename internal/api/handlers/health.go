package handlers

import (
	"net/http"
	"strade/internal/constants"
)

func (m *HandlerManager) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	msg := m.Translator.GetMessageOrDefault(r.Context(), constants.HealthOk, "Ok", nil)

	m.Logger.Info(msg)

	if err := m.Transporter.JSONResponse(w, http.StatusOK, map[string]any{"message": msg}); err != nil {
		m.Transporter.InternalServerErrorBasic(w, r, err)
	}
}
