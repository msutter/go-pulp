package pulp

import ()

type PublishTask struct {
	Id             string          `json:"task_id"`
	StartTime      string          `json:"start_time"`
	FinishTime     string          `json:"finish_time"`
	State          string          `json:"state"`
	Error          *PublishError   `json:"error"`
	ProgressReport *ProgressReport `json:"progress_report"`
	Result         *PublishResult  `json:"result"`
}

type PublishError struct {
	Description string `json:"description"`
}

type PublishResult struct {
	Details []*Distributor `json:"details"`
}

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
