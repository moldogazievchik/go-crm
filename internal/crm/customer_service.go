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

var ErrCustomerHasLeads = ErrValidation("customer has leads and cannot be deleted")

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *CustomerService) UpdateCustomer(id, name, email string) (Customer, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Customer{}, ErrValidation("id is required")
	}

	current, err := s.repo.GetByID(id)
	if err != nil {
		return Customer{}, err
	}

	// частичное обновление: если поле пустое — не трогаем
	if strings.TrimSpace(name) != "" {
		current.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(email) != "" {
		e := strings.TrimSpace(email)
		if !strings.Contains(e, "@") {
			return Customer{}, ErrValidation("email is invalid")
		}
		current.Email = e
	}

	return s.repo.Update(current)
}

func (s *CustomerService) DeleteCustomer(id string, leadRepo LeadRepository) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrValidation("id is required")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	leads, err := leadRepo.ListByCustomerID(id)
	if err != nil {
		return err
	}
	if len(leads) > 0 {
		return ErrCustomerHasLeads
	}

	return s.repo.Delete(id)
}
