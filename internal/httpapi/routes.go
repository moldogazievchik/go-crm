package httpapi

import (
	"net/http"
	"strings"

	"github.com/moldogazievchik/go-crm/internal/crm"
)

// Routes собирает все маршруты приложения.
func Routes() http.Handler {
	mux := http.NewServeMux()

	// Health endpoint
	mux.HandleFunc("/health", healthHandler)

	// ---------- FILE STORE ----------
	store := crm.NewStore("./data.json")
	_ = store.Load()

	// ---------- Customers ----------
	customerRepo := crm.NewFileCustomerRepo(store)
	customerSvc := crm.NewCustomerService(customerRepo)

	// ---------- Leads ----------
	leadRepo := crm.NewFileLeadRepo(store)
	leadSvc := crm.NewLeadService(leadRepo, customerRepo)

	// ---------- Handlers ----------
	ch := NewCustomerHandler(customerSvc, leadSvc)
	lh := NewLeadHandler(leadSvc)

	mux.HandleFunc("/customers", ch.customers)
	//mux.HandleFunc("/customers/", ch.customerByID)
	//mux.HandleFunc("/customers/", ch.customerLeads)

	mux.HandleFunc("/customers/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/leads") {
			ch.customerLeads(w, r)
			return
		}
		ch.customerByID(w, r)
	})

	mux.HandleFunc("/leads", lh.leads)
	mux.HandleFunc("/leads/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/status") {
			lh.leadStatus(w, r)
			return
		}
		lh.leadByID(w, r)
	})

	return mux
}
