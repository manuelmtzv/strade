package handle

import (
	"net/http"
	"strade/internal/api/transport"
	"strade/internal/constants"
	"strade/internal/models"
)

// @Summary Get settlement types
// @Description Get a list of settlement types
// @Tags SettlementTypes
// @Accept json
// @Produce json
// @Success 200 {object} models.SettlementType
// @Failure 500 {object} transport.ErrorResponse
// @Router /settlementtypes [get]
func (h *Handler) HandleGetSettlementTypes(w http.ResponseWriter, r *http.Request) {
	settlementTypes, err := h.Store.SettlementTypeStore.FindAll()
	if err != nil {
		h.Transporter.InternalServerError(w, r,
			transport.WithTechnicalError(err),
			transport.WithMessageID(constants.ErrorInternalServerError),
			transport.WithErrorCode(constants.ErrCodeInternalServerError),
		)
		return
	}

	metadata := models.ListMetadata{
		Total: len(settlementTypes),
	}

	if err := h.Transporter.JSONResponseWithMetadata(w, http.StatusOK, settlementTypes, metadata); err != nil {
		h.Transporter.InternalServerError(w, r,
			transport.WithTechnicalError(err),
			transport.WithMessageID(constants.ErrorInternalServerError),
			transport.WithErrorCode(constants.ErrCodeInternalServerError),
		)
	}
}
