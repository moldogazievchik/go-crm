package crm

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

// CustomerService — бизнес-логика по клиентам.
type CustomerService struct {
	repo CustomerRepository
}

func NewCustomerService(repo CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) CreateCustomer(name, email string) (Customer, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return Customer{}, ErrValidation("name is required")
	}

	if email == "" || !strings.Contains(email, "@") {
		return Customer{}, ErrValidation("email is invalid")
	}

	c := Customer{
		ID:        newID(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	return s.repo.Create(c)
}

func (s *CustomerService) ListCustomers() ([]Customer, error) {
	return s.repo.List()
}

func (s *CustomerService) GetCustomer(id string) (Customer, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Customer{}, ErrValidation("id is required")
	}
	return s.repo.GetByID(id)
}

type ErrValidation string

func (e ErrValidation) Error() string { return string(e) }

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
