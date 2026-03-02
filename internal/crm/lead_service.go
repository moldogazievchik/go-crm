package crm

import (
	"strings"
	"time"
)

type LeadService struct {
	leads     LeadRepository
	customers CustomerRepository
}

func NewLeadService(leads LeadRepository, customers CustomerRepository) *LeadService {
	return &LeadService{leads: leads, customers: customers}
}

func (s *LeadService) CreateLead(customerID, title string, value int) (Lead, error) {
	customerID = strings.TrimSpace(customerID)
	title = strings.TrimSpace(title)

	if customerID == "" {
		return Lead{}, ErrValidation("customer_id is required")
	}
	if title == "" {
		return Lead{}, ErrValidation("title is required")
	}
	if value < 0 {
		return Lead{}, ErrValidation("value must be >= 0")
	}

	// Проверяем, что customer существует
	_, err := s.customers.GetByID(customerID)
	if err != nil {
		if err == ErrCustomerNotFound {
			return Lead{}, ErrValidation("customer_id does not exist")
		}
		return Lead{}, err
	}

	l := Lead{
		ID:         newID(),
		CustomerID: customerID,
		Title:      title,
		Status:     LeadNew,
		Value:      value,
		CreatedAt:  time.Now(),
	}

	return s.leads.Create(l)
}

func (s *LeadService) ListLeads() ([]Lead, error) {
	return s.leads.List()
}

func (s *LeadService) ListLeadsByCustomer(customerID string) ([]Lead, error) {
	customerID = strings.TrimSpace(customerID)
	if customerID == "" {
		return nil, ErrValidation("customer_id is required")
	}

	_, err := s.customers.GetByID(customerID)
	if err != nil {
		if err == ErrCustomerNotFound {
			return nil, ErrValidation("customer_id does not exist")
		}
		return nil, err
	}

	return s.leads.ListByCustomerID(customerID)
}

func (s *LeadService) GetLead(id string) (Lead, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Lead{}, ErrValidation("id is required")
	}
	return s.leads.GetByID(id)
}

func (s *LeadService) UpdateStatus(id string, status LeadStatus) (Lead, error) {
	switch status {
	case LeadNew, LeadInProgress, LeadWon, LeadLost:
		// ok
	default:
		return Lead{}, ErrValidation("invalid status")
	}

	l, err := s.leads.GetByID(id)
	if err != nil {
		return Lead{}, err
	}

	l.Status = status
	return s.leads.Update(l)
}

func (s *LeadService) UpdateLead(id string, title *string, value *int, customerID *string) (Lead, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Lead{}, ErrValidation("id is required")
	}

	current, err := s.leads.GetByID(id)
	if err != nil {
		return Lead{}, err
	}

	// title (если передали)
	if title != nil {
		t := strings.TrimSpace(*title)
		if t == "" {
			return Lead{}, ErrValidation("title is required")
		}
		current.Title = t
	}

	// value (если передали)
	if value != nil {
		if *value < 0 {
			return Lead{}, ErrValidation("value must be >= 0")
		}
		current.Value = *value
	}

	// customer_id (если передали) + проверка что клиент существует
	if customerID != nil {
		cid := strings.TrimSpace(*customerID)
		if cid == "" {
			return Lead{}, ErrValidation("customer_id is required")
		}

		_, err := s.customers.GetByID(cid)
		if err != nil {
			if err == ErrCustomerNotFound {
				return Lead{}, ErrValidation("customer_id does not exist")
			}
			return Lead{}, err
		}

		current.CustomerID = cid
	}

	return s.leads.Update(current)
}
