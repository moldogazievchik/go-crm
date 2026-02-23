package crm

type FileLeadRepo struct {
	store *Store
}

func NewFileLeadRepo(store *Store) *FileLeadRepo {
	return &FileLeadRepo{store: store}
}

func (r *FileLeadRepo) Create(l Lead) (Lead, error) {
	r.store.mu.Lock()
	r.store.Leads[l.ID] = l
	r.store.mu.Unlock()

	if err := r.store.Save(); err != nil {
		return Lead{}, err
	}

	return l, nil
}

func (r *FileLeadRepo) List() ([]Lead, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	out := make([]Lead, 0, len(r.store.Leads))
	for _, l := range r.store.Leads {
		out = append(out, l)
	}
	return out, nil
}

func (r *FileLeadRepo) ListByCustomerID(customerID string) ([]Lead, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	out := make([]Lead, 0)
	for _, l := range r.store.Leads {
		if l.CustomerID == customerID {
			out = append(out, l)
		}
	}
	return out, nil
}

func (r *FileLeadRepo) GetByID(id string) (Lead, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	l, ok := r.store.Leads[id]
	if !ok {
		return Lead{}, ErrLeadNotFound
	}
	return l, nil
}

func (r *FileLeadRepo) Update(l Lead) (Lead, error) {
	r.store.mu.Lock()
	if _, ok := r.store.Leads[l.ID]; !ok {
		r.store.mu.Unlock()
		return Lead{}, ErrLeadNotFound
	}
	r.store.Leads[l.ID] = l
	r.store.mu.Unlock()

	if err := r.store.Save(); err != nil {
		return Lead{}, err
	}
	return l, nil
}
