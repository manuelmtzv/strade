package handle

import (
	"fmt"
	"net/http"
	"strade/internal/api/transport"
	"strade/internal/constants"
	"strade/internal/models"
	"strade/internal/utils"
)

func (h *Handler) HandleSearchSettlements(w http.ResponseWriter, r *http.Request) {
	query, _ := h.Transporter.GetSearchParam(r, "q")
	if query == nil || *query == "" {
		h.Transporter.BadRequest(w, r,
			transport.WithMessageID(constants.ErrorInvalidQuery),
			transport.WithErrorCode(constants.ErrCodeInvalidQuery),
		)
		return
	}

	var limit *int
	if limit, _ = h.Transporter.GetSearchParamInt(r, "limit"); limit != nil {
		if *limit < 1 {
			limit = utils.IntPointer(10)
		}
	} else {
		limit = utils.IntPointer(10)
	}

	settlements, err := h.Store.SettlementStore.SearchByName(r.Context(), *query, *limit)
	if err != nil {
		h.Transporter.InternalServerError(w, r,
			transport.WithTechnicalError(fmt.Errorf("error searching settlements: %w", err)),
			transport.WithMessageID(constants.ErrorInternalServerError),
			transport.WithErrorCode(constants.ErrCodeInternalServerError),
		)
		return
	}

	results := make([]models.SettlementSearchResult, len(settlements))
	for i, s := range settlements {
		results[i] = models.SettlementSearchResult{
			Name:           s.Name,
			PostalCode:     s.PostalCode,
			SettlementType: s.SettlementType.Name,
			Municipality:   s.Municipality.Name,
			State:          s.State.Name,
			City:           s.City.Name,
		}
	}

	metadata := models.ListMetadata{
		Total: len(settlements),
		Limit: *limit,
	}

	if err := h.Transporter.JSONResponseWithMetadata(w, http.StatusOK, results, metadata); err != nil {
		h.Transporter.InternalServerError(w, r,
			transport.WithTechnicalError(err),
			transport.WithMessageID(constants.ErrorInternalServerError),
			transport.WithErrorCode(constants.ErrCodeInternalServerError),
		)
	}
}
