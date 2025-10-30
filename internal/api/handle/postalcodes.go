package handle

import (
	"database/sql"
	"fmt"
	"net/http"
	"strade/internal/api/transport"
	"strade/internal/constants"
	"strade/internal/models"
)

func (h *Handler) HandleGetPostalCodeSettlements(w http.ResponseWriter, r *http.Request) {
	postalCode, _ := h.Transporter.GetUrlParam(r, "postalCode")

	if postalCode == nil || *postalCode == "" {
		h.Transporter.BadRequest(w, r,
			transport.WithMessageID(constants.ErrorInvalidPostalCode),
			transport.WithErrorCode(constants.ErrCodeInvalidPostalCode),
		)
		return
	}

	settlements, err := h.Store.SettlementStore.FindByPostalCode(r.Context(), *postalCode)
	if err != nil && err != sql.ErrNoRows {
		h.Transporter.InternalServerError(w, r,
			transport.WithTechnicalError(fmt.Errorf("error getting settlement by postal code %s: %w", *postalCode, err)),
			transport.WithMessageID(constants.ErrorInternalServerError),
			transport.WithErrorCode(constants.ErrCodeInternalServerError),
		)
		return
	}

	if len(settlements) == 0 {
		h.Transporter.NotFound(w, r,
			transport.WithTechnicalError(fmt.Errorf("postal code not found: %s", *postalCode)),
			transport.WithMessageID(constants.ErrorPostalCodeNotFound),
			transport.WithErrorCode(constants.ErrCodePostalCodeNotFound),
		)
		return
	}

	postalCodeDetails := make([]models.PostalCodeDetails, len(settlements))
	for i, settlement := range settlements {
		postalCodeDetails[i] = models.PostalCodeDetails{
			PostalCode:     *postalCode,
			Settlement:     settlement.Name,
			SettlementType: settlement.SettlementType.Name,
			Municipality:   settlement.Municipality.Name,
			State:          settlement.State.Name,
			City:           settlement.City.Name,
		}
	}

	metadata := models.ListMetadata{
		Total: len(settlements),
	}

	if err := h.Transporter.JSONResponseWithMetadata(w, http.StatusOK, postalCodeDetails, metadata); err != nil {
		h.Transporter.InternalServerError(w, r,
			transport.WithTechnicalError(err),
			transport.WithMessageID(constants.ErrorInternalServerError),
			transport.WithErrorCode(constants.ErrCodeInternalServerError),
		)
	}
}
