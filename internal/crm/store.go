package crm

import (
	"encoding/json"
	"os"
	"sync"
)

// Store — общий контейнер данных (и customers, и leads).
type Store struct {
	mu sync.RWMutex

	path string

	Customers map[string]Customer `json:"customers"`
	Leads     map[string]Lead     `json:"leads"`
}

func NewStore(path string) *Store {
	return &Store{
		path:      path,
		Customers: make(map[string]Customer),
		Leads:     make(map[string]Lead),
	}
}

// Load читает данные из файла, если он существует.
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			// файла нет — это нормально, стартуем с пустыми данными
			return nil
		}
		return err
	}

	// Декодим JSON прямо в структуру
	if err := json.Unmarshal(b, s); err != nil {
		return err
	}

	// на всякий случай, если в файле пусто
	if s.Customers == nil {
		s.Customers = make(map[string]Customer)
	}
	if s.Leads == nil {
		s.Leads = make(map[string]Lead)
	}

	return nil
}

// Save сохраняет текущее состояние в файл (атомарно через temp + rename).
func (s *Store) Save() error {
	// Сначала делаем снапшот под RLock
	s.mu.RLock()
	snapshot := struct {
		Customers map[string]Customer `json:"customers"`
		Leads     map[string]Lead     `json:"leads"`
	}{
		Customers: copyCustomers(s.Customers),
		Leads:     copyLeads(s.Leads),
	}
	s.mu.RUnlock()

	b, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func copyCustomers(src map[string]Customer) map[string]Customer {
	dst := make(map[string]Customer, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func copyLeads(src map[string]Lead) map[string]Lead {
	dst := make(map[string]Lead, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
