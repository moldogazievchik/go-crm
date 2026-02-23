package crm

import (
	"errors"
	"sync"
)

type LeadRepository interface {
	Create(l Lead) (Lead, error)
	List() ([]Lead, error)
	ListByCustomerID(customerID string) ([]Lead, error)
	GetByID(id string) (Lead, error)
	Update(l Lead) (Lead, error)
}

var ErrLeadNotFound = errors.New("lead not found")

type MemoryLeadRepo struct {
	mu   sync.RWMutex
	byID map[string]Lead
}

func NewMemoryLeadRepo() *MemoryLeadRepo {
	return &MemoryLeadRepo{
		byID: make(map[string]Lead),
	}
}

func (r *MemoryLeadRepo) Create(l Lead) (Lead, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[l.ID] = l
	return l, nil
}

func (r *MemoryLeadRepo) List() ([]Lead, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Lead, 0, len(r.byID))
	for _, l := range r.byID {
		out = append(out, l)
	}
	return out, nil
}

func (r *MemoryLeadRepo) ListByCustomerID(customerID string) ([]Lead, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Lead, 0)
	for _, l := range r.byID {
		if l.CustomerID == customerID {
			out = append(out, l)
		}
	}
	return out, nil
}

func (r *MemoryLeadRepo) GetByID(id string) (Lead, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	l, ok := r.byID[id]
	if !ok {
		return Lead{}, ErrLeadNotFound
	}
	return l, nil
}

func (r *MemoryLeadRepo) Update(l Lead) (Lead, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[l.ID]; !ok {
		return Lead{}, ErrLeadNotFound
	}

	r.byID[l.ID] = l
	return l, nil
}
