package handle

import (
	"database/sql"
	"fmt"
	"net/http"
	"strade/internal/api/transport"
	"strade/internal/constants"
)

func (h *Handler) HandleGetPostalCodeSettlements(w http.ResponseWriter, r *http.Request) {
	postalCode, err := h.Transporter.GetUrlParam(r, "postalCode")
	if err != nil {
		h.Transporter.InternalServerError(w, r,
			transport.WithTechnicalError(fmt.Errorf("error getting url param: %w", err)),
			transport.WithMessageID(constants.ErrorInternalServerError),
			transport.WithErrorCode(constants.ErrCodeInternalServerError),
		)
		return
	}

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

	if err := h.Transporter.JSONResponse(w, http.StatusOK, settlements); err != nil {
		h.Transporter.InternalServerError(w, r,
			transport.WithTechnicalError(err),
			transport.WithMessageID(constants.ErrorInternalServerError),
			transport.WithErrorCode(constants.ErrCodeInternalServerError),
		)
	}
}
