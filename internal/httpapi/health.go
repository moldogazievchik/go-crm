package httpapi

import (
	"encoding/json"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Гооврим клиенту, чо ответ будет JSON
	w.Header().Set("Content-Type", "application/json")

	// Пишем статус 200 OK
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
