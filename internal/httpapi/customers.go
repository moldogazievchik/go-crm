package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/moldogazievchik/go-crm/internal/crm"
)

type CustomerHandler struct {
	svc      *crm.CustomerService
	leadSvc  *crm.LeadService
	leadRepo crm.LeadRepository
}

func NewCustomerHandler(svc *crm.CustomerService, leadSvc *crm.LeadService, leadRepo crm.LeadRepository) *CustomerHandler {
	return &CustomerHandler{
		svc:      svc,
		leadSvc:  leadSvc,
		leadRepo: leadRepo}
}

func (h *CustomerHandler) customers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createCustomer(w, r)
	case http.MethodGet:
		h.listCustomers(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *CustomerHandler) customerByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimSpace(chi.URLParam(r, "id"))

	c, err := h.svc.GetCustomer(id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, c)
}

type createCustomerrequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	var req createCustomerrequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	c, err := h.svc.CreateCustomer(req.Name, req.Email)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, c)
}

func (h *CustomerHandler) listCustomers(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListCustomers()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	var vErr crm.ErrValidation
	if errors.As(err, &vErr) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": vErr.Error()})
		return
	}

	if errors.Is(err, crm.ErrCustomerNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
}

func (h *CustomerHandler) customerLeads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// /customers/{id}/leads
	path := strings.TrimPrefix(r.URL.Path, "/customers/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "leads" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	customerID := chi.URLParam(r, "id")

	items, err := h.leadSvc.ListLeadsByCustomer(customerID)
	if err != nil {
		writeDomainError(w, err) // твой writeError уже умеет ErrValidation/NotFound
		return
	}

	writeJSON(w, http.StatusOK, items)
}

type patchCustomerRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func (h *CustomerHandler) patchCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")

	var req patchCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_json", "invalid json")
		return
	}

	name := ""
	email := ""
	if req.Name != nil {
		name = *req.Name
	}
	if req.Email != nil {
		email = *req.Email
	}

	updated, err := h.svc.UpdateCustomer(id, name, email)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *CustomerHandler) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")

	err := h.svc.DeleteCustomer(id, h.leadRepo)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
