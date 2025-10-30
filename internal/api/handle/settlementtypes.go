package handle

import (
	"net/http"
	"strade/internal/api/transport"
	"strade/internal/constants"
	"strade/internal/models"
)

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
