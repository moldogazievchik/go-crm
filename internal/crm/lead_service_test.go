package crm

import "testing"

func TestLeadService_CreateLead_CustomerMustExist(t *testing.T) {
	customerRepo := NewMemoryCustomerRepo()
	leadRepo := NewMemoryLeadRepo()

	leadSvc := NewLeadService(leadRepo, customerRepo)

	_, err := leadSvc.CreateLead("missing-customer", "Deal", 100)
	if err == nil {
		t.Fatalf("expected error when customer does not exist")
	}
}

func TestLeadService_CreateUpdateLead(t *testing.T) {
	customerRepo := NewMemoryCustomerRepo()
	leadRepo := NewMemoryLeadRepo()

	customerSvc := NewCustomerService(customerRepo)
	leadSvc := NewLeadService(leadRepo, customerRepo)

	c, err := customerSvc.CreateCustomer("Aktan", "aktan@example.com")
	if err != nil {
		t.Fatalf("create customer error: %v", err)
	}

	l, err := leadSvc.CreateLead(c.ID, "Website redesign", 500)
	if err != nil {
		t.Fatalf("create lead error: %v", err)
	}
	if l.Status != LeadNew {
		t.Fatalf("expected status %q got %q", LeadNew, l.Status)
	}

	// UpdateLead: меняем title и value
	newTitle := "Website redesign v2"
	newValue := 700
	updated, err := leadSvc.UpdateLead(l.ID, &newTitle, &newValue, nil)
	if err != nil {
		t.Fatalf("update lead error: %v", err)
	}
	if updated.Title != newTitle {
		t.Fatalf("title not updated")
	}
	if updated.Value != newValue {
		t.Fatalf("value not updated")
	}

	// UpdateStatus
	updated2, err := leadSvc.UpdateStatus(l.ID, LeadInProgress)
	if err != nil {
		t.Fatalf("update status error: %v", err)
	}
	if updated2.Status != LeadInProgress {
		t.Fatalf("expected status %q got %q", LeadInProgress, updated2.Status)
	}
}
