package crm

type FileCustomerRepo struct {
	store *Store
}

func NewFileCustomerRepo(store *Store) *FileCustomerRepo {
	return &FileCustomerRepo{store: store}
}

func (r *FileCustomerRepo) Create(c Customer) (Customer, error) {
	r.store.mu.Lock()
	r.store.Customers[c.ID] = c
	r.store.mu.Unlock()

	// сохраняем на диск
	if err := r.store.Save(); err != nil {
		return Customer{}, err
	}

	return c, nil
}

func (r *FileCustomerRepo) List() ([]Customer, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	out := make([]Customer, 0, len(r.store.Customers))
	for _, c := range r.store.Customers {
		out = append(out, c)
	}
	return out, nil
}

func (r *FileCustomerRepo) GetByID(id string) (Customer, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	c, ok := r.store.Customers[id]
	if !ok {
		return Customer{}, ErrCustomerNotFound
	}
	return c, nil
}

func (r *FileCustomerRepo) Update(c Customer) (Customer, error) {
	r.store.mu.Lock()
	if _, ok := r.store.Customers[c.ID]; !ok {
		r.store.mu.Unlock()
		return Customer{}, ErrCustomerNotFound
	}
	r.store.Customers[c.ID] = c
	r.store.mu.Unlock()

	if err := r.store.Save(); err != nil {
		return Customer{}, err
	}
	return c, nil
}
