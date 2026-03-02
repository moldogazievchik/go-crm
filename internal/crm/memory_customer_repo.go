package crm

import (
	"errors"
	"sync"
)

// CustomerRepository - интерфейс хранилица клиентов.
type CustomerRepository interface {
	Create(c Customer) (Customer, error)
	List() ([]Customer, error)
	GetByID(id string) (Customer, error)
	Update(c Customer) (Customer, error)
}

var ErrCustomerNotFound = errors.New("customer not found")

type MemoryCustomerRepo struct {
	mu   sync.RWMutex
	byID map[string]Customer
}

// NewMemoryCustomerRepo — конструктор.
func NewMemoryCustomerRepo() *MemoryCustomerRepo {
	return &MemoryCustomerRepo{
		byID: make(map[string]Customer),
	}
}

func (r *MemoryCustomerRepo) Create(c Customer) (Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[c.ID] = c
	return c, nil
}

func (r *MemoryCustomerRepo) List() ([]Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Customer, 0, len(r.byID))
	for _, c := range r.byID {
		out = append(out, c)
	}
	return out, nil
}
func (r *MemoryCustomerRepo) GetByID(id string) (Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.byID[id]
	if !ok {
		return Customer{}, ErrCustomerNotFound
	}
	return c, nil
}

func (r *MemoryCustomerRepo) Update(c Customer) (Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[c.ID]; !ok {
		return Customer{}, ErrCustomerNotFound
	}
	r.byID[c.ID] = c
	return c, nil
}
