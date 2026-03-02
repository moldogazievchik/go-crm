package httpapi

import (
	"errors"
	"net/http"

	"github.com/moldogazievchik/go-crm/internal/crm"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiErrorResponse struct {
	Error apiError `json:"error"`
}

func writeAPIError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, apiErrorResponse{
		Error: apiError{Code: code, Message: msg},
	})
}

func writeDomainError(w http.ResponseWriter, err error) {
	var vErr crm.ErrValidation
	if errors.As(err, &vErr) {
		writeAPIError(w, http.StatusBadRequest, "validation_error", vErr.Error())
		return
	}

	if errors.Is(err, crm.ErrCustomerNotFound) || errors.Is(err, crm.ErrLeadNotFound) {
		writeAPIError(w, http.StatusNotFound, "not_found", "not found")
		return
	}

	writeAPIError(w, http.StatusInternalServerError, "internal_error", "internal")
}
