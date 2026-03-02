package crm

import "testing"

func TestCustomerService_CreateCustomer_Validation(t *testing.T) {
	repo := NewMemoryCustomerRepo()
	svc := NewCustomerService(repo)

	_, err := svc.CreateCustomer("", "a@b.com")
	if err == nil {
		t.Fatalf("expected error for empty name")
	}

	_, err = svc.CreateCustomer("Aktan", "not-an-email")
	if err == nil {
		t.Fatalf("expected error for invalid email")
	}
}

func TestCustomerService_CreateAndGet(t *testing.T) {
	repo := NewMemoryCustomerRepo()
	svc := NewCustomerService(repo)

	created, err := svc.CreateCustomer("Aktan", "aktan@example.com")
	if err != nil {
		t.Fatalf("create error: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected id to be set")
	}

	got, err := svc.GetCustomer(created.ID)
	if err != nil {
		t.Fatalf("get error: %v", err)
	}
	if got.Email != "aktan@example.com" {
		t.Fatalf("expected email %q, got %q", "aktan@example.com", got.Email)
	}
}

func TestCustomerService_UpdateCustomer_Partial(t *testing.T) {
	repo := NewMemoryCustomerRepo()
	svc := NewCustomerService(repo)

	created, err := svc.CreateCustomer("Aktan", "aktan@example.com")
	if err != nil {
		t.Fatalf("create error: %v", err)
	}

	// меняем только name
	updated, err := svc.UpdateCustomer(created.ID, "Aktan Updated", "")
	if err != nil {
		t.Fatalf("update error: %v", err)
	}
	if updated.Name != "Aktan Updated" {
		t.Fatalf("expected updated name, got %q", updated.Name)
	}
	if updated.Email != "aktan@example.com" {
		t.Fatalf("email should remain unchanged, got %q", updated.Email)
	}
}
