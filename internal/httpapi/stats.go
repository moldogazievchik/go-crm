package httpapi

import (
	"net/http"

	"github.com/moldogazievchik/go-crm/internal/crm"
)

type StatsHandler struct {
	leadSvc *crm.LeadService
}

func NewStatsHandler(leadSvc *crm.LeadService) *StatsHandler {
	return &StatsHandler{leadSvc: leadSvc}
}

func (h *StatsHandler) stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	result, err := h.leadSvc.Stats()
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "internal_error", "internal")
		return
	}

	writeJSON(w, http.StatusOK, result)
}
