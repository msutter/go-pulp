//
// Copyright 2016, Marc Sutter
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package pulp

import (
	"fmt"
)

type TasksService struct {
	client *Client
}

func (s *TasksService) ListTasks() ([]*Task, *Response, error) {
	req, err := s.client.NewRequest("GET", "tasks/", nil)
	if err != nil {
		return nil, nil, err
	}

	var t []*Task
	resp, err := s.client.Do(req, &t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}

func (s *TasksService) GetTask(task string) (*Task, *Response, error) {
	u := fmt.Sprintf("tasks/%s/", task)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	t := new(Task)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}

type Task struct {
	Id             string          `json:"task_id"`
	StartTime      string          `json:"start_time"`
	FinishTime     string          `json:"finish_time"`
	State          string          `json:"state"`
	Error          *Error          `json:"error"`
	ProgressReport *ProgressReport `json:"progress_report"`
	Result         *Result         `json:"result"`
}

func (t *Task) String() string {
	return Stringify(t)
}

func (t *Task) Importer() (importer string) {
	if t.ProgressReport.YumImporter != nil {
		importer = "yum"
	}

	if t.ProgressReport.DockerImporter != nil {
		importer = "docker"
	}
	return
}

type ProgressReport struct {
	// yum importer
	YumImporter *Importer `json:"yum_importer"`
	// docker importer
	DockerImporter *Importer `json:"docker_importer"`
}

type Details struct {
	Content *Content `json:"content"`
}

type Result struct {
	Details     `json:"details"`
	ResultUnits []ResultUnit `json:"units_successful"`
}

type ResultUnit struct {
	UnitKey *UnitKey `json:"unit_key"`
	TypeId  string   `json:"type_id"`
}

type UnitKey struct {
	Id       string `json:"id"`
	FileName string `json:"filename"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Release  string `json:"release"`
	Arch     string `json:"arch"`
}
