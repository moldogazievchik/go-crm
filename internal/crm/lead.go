package crm

import "time"

type LeadStatus string

const (
	LeadNew        LeadStatus = "new"
	LeadInProgress LeadStatus = "in_progress"
	LeadWon        LeadStatus = "won"
	LeadLost       LeadStatus = "lost"
)

type Lead struct {
	ID         string     `json:"id"`
	CustomerID string     `json:"customer_id"`
	Title      string     `json:"title"`
	Status     LeadStatus `json:"status"`
	Value      int        `json:"value"`
	CreatedAt  time.Time  `json:"created_at"`
}
