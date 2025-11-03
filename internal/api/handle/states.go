package handle

import (
	"net/http"
	"strade/internal/api/transport"
	"strade/internal/models"

	"github.com/go-chi/chi/v5"
)

// @Summary Get all states
// @Description Get a list of all states
// @Tags States
// @Accept json
// @Produce json
// @Success 200 {object} models.StatesResponse
// @Failure 500 {object} transport.ErrorResponse
// @Router /states [get]
func (h *Handler) HandleGetStates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	states, err := h.Store.StateStore.FindAll(ctx)
	if err != nil {
		h.Transporter.SendError(w, r, err)
		return
	}

	if states == nil {
		states = []*models.State{} // Return empty array instead of null
	}

	h.Transporter.SendJSON(w, r, http.StatusOK, models.StatesResponse{Data: states})
}

// @Summary Get municipalities by state
// @Description Get a list of municipalities for a specific state
// @Tags States
// @Accept json
// @Produce json
// @Param stateId path string true "State ID"
// @Success 200 {object} models.MunicipalitiesResponse
// @Failure 400 {object} transport.ErrorResponse
// @Failure 404 {object} transport.ErrorResponse
// @Failure 500 {object} transport.ErrorResponse
// @Router /states/{stateId}/municipalities [get]
func (h *Handler) HandleGetStateMunicipalities(w http.ResponseWriter, r *http.Request) {
	stateID := chi.URLParam(r, "stateId")
	if stateID == "" {
		h.Transporter.BadRequest(w, r, transport.WithMessage("State ID is required"))
		return
	}

	municipalities, err := h.Store.StateStore.FindMunicipalitiesByStateID(r.Context(), stateID)
	if err != nil {
		h.Transporter.InternalServerError(w, r, transport.WithMessage("Failed to fetch municipalities"), transport.WithTechnicalError(err))
		return
	}

	if municipalities == nil {
		municipalities = []*models.Municipality{}
	}

	h.Transporter.SendJSON(w, r, http.StatusOK, models.MunicipalitiesResponse{Data: municipalities})
}
