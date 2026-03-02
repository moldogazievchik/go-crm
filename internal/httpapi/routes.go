package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/moldogazievchik/go-crm/internal/crm"
)

// Routes собирает все маршруты приложения.
func Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health endpoint
	r.Get("/health", healthHandler)

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

	// Customers
	r.Route("/customers", func(r chi.Router) {
		r.Get("/", ch.customers)
		r.Post("/", ch.customers)
		r.Get("/{id}", ch.customerByID)
		r.Patch("/{id}", ch.patchCustomer)
		r.Get("/{id}/leads", ch.customerLeads)
	})

	// Leads
	r.Route("/leads", func(r chi.Router) {
		r.Get("/", lh.leads)
		r.Post("/", lh.leads)
		r.Get("/{id}", lh.leadByID)
		r.Patch("/{id}", lh.patchLead)
		r.Patch("/{id}/status", lh.leadStatus)
	})

	return r
}
