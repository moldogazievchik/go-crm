package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/moldogazievchik/go-crm/internal/crm"
)

type LeadHandler struct {
	svc *crm.LeadService
}

func NewLeadHandler(svc *crm.LeadService) *LeadHandler {
	return &LeadHandler{svc: svc}
}

// POST /leads, GET /leads
func (h *LeadHandler) leads(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createLead(w, r)
	case http.MethodGet:
		h.listLeads(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GET /leads/{id}
func (h *LeadHandler) leadByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/leads/")
	id = strings.TrimSpace(id)

	l, err := h.svc.GetLead(id)
	if err != nil {
		writeLeadError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, l)
}

// PATCH /leads/{id}/status
func (h *LeadHandler) leadStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// /leads/{id}/status
	path := strings.TrimPrefix(r.URL.Path, "/leads/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "status" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	id := strings.TrimSpace(parts[0])

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	l, err := h.svc.UpdateStatus(id, crm.LeadStatus(req.Status))
	if err != nil {
		writeLeadError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, l)
}

type createLeadRequest struct {
	CustomerID string `json:"customer_id"`
	Title      string `json:"title"`
	Value      int    `json:"value"`
}

func (h *LeadHandler) createLead(w http.ResponseWriter, r *http.Request) {
	var req createLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	l, err := h.svc.CreateLead(req.CustomerID, req.Title, req.Value)
	if err != nil {
		writeLeadError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, l)
}

func (h *LeadHandler) listLeads(w http.ResponseWriter, r *http.Request) {
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	customerID := strings.TrimSpace(r.URL.Query().Get("customer_id"))

	items, err := h.svc.ListLeads()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
		return
	}

	// фильтрация в HTTP слое (позже перенесём в сервис/репо)
	filtered := make([]crm.Lead, 0, len(items))
	for _, l := range items {
		if status != "" && string(l.Status) != status {
			continue
		}
		if customerID != "" && l.CustomerID != customerID {
			continue
		}
		filtered = append(filtered, l)
	}

	writeJSON(w, http.StatusOK, filtered)
}

func writeLeadError(w http.ResponseWriter, err error) {
	var vErr crm.ErrValidation
	if errors.As(err, &vErr) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": vErr.Error()})
		return
	}
	if errors.Is(err, crm.ErrLeadNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
}
