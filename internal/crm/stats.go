package crm

type LeadStats struct {
	NewCount        int `json:"new_count"`
	InProgressCount int `json:"in_progress_count"`
	WonCount        int `json:"won_count"`
	LostCount       int `json:"lost_count"`
	TotalCount      int `json:"total_count"`
	TotalValue      int `json:"total_value"`
	WonValue        int `json:"won_value"`
}
