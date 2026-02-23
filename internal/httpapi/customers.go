package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/moldogazievchik/go-crm/internal/crm"
)

type CustomerHandler struct {
	svc     *crm.CustomerService
	leadSvc *crm.LeadService
}

func NewCustomerHandler(svc *crm.CustomerService, leadSvc *crm.LeadService) *CustomerHandler {
	return &CustomerHandler{svc: svc, leadSvc: leadSvc}

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

	id := strings.TrimPrefix(r.URL.Path, "/customers/")
	id = strings.TrimSpace(id)

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

	customerID := strings.TrimSpace(parts[0])

	items, err := h.leadSvc.ListLeadsByCustomer(customerID)
	if err != nil {
		writeError(w, err) // твой writeError уже умеет ErrValidation/NotFound
		return
	}

	writeJSON(w, http.StatusOK, items)
}
