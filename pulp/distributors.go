package pulp

import ()

type Distributor struct {
	Description  string   `json:"description"`
	ItemsTotal   int      `json:"items_total"`
	NumSuccess   int      `json:"num_success"`
	NumFailures  int      `json:"num_failures"`
	NumProcessed int      `json:"num_processed"`
	StepType     string   `json:"step_type"`
	StepId       string   `json:"step_id"`
	State        string   `json:"state"`
	ErrorDetails []string `json:"error_details"`
}
