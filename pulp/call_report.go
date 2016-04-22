package pulp

import ()

// Pulp Api docs:
// http://pulp.readthedocs.org/en/latest/dev-guide/conventions/sync-v-async.html#call-report
type CallReport struct {
	Result       string `json:"result"`
	Error        *Error `json:"error"`
	SpawnedTasks []struct {
		Href   string `json:"_href"`
		TaskId string `json:"task_id"`
	} `json:"spawned_tasks"`
}
