package pulp

import ()

type Importer struct {
	Id             string          `json:"id"`
	ImporterConfig *ImporterConfig `json:"config"`
	Content        *Content        `json:"content"`
	Metadata       *Metadata       `json:"metadata"`
}

type ImporterConfig struct {
	Feed          string `json:"feed"`
	RemoveMissing bool   `json:"remove_missing"`
}

// included in task
type Content struct {
	State        string   `json:"state"`
	ItemsTotal   int      `json:"items_total"`
	ItemsLeft    int      `json:"items_left"`
	SizeTotal    int      `json:"size_total"`
	SizeLeft     int      `json:"size_left"`
	ErrorDetails []string `json:"error_details"`
}

// included in task
type Metadata struct {
	State string
	Error string
}
