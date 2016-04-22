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
	"encoding/json"
	"fmt"
)

type RepositoriesService struct {
	client *Client
}

type Repository struct {
	Id        string      `json:"id"`
	Name      string      `json:"display_name"`
	Importers []*Importer `json:"importers"`
}

func (r Repository) String() string {
	return Stringify(r)
}

func (s *RepositoriesService) ListRepositories(opt *GetRepositoryOptions) ([]*Repository, *Response, error) {

	req, err := s.client.NewRequest("GET", "repositories/", opt)
	if err != nil {
		return nil, nil, err
	}

	var r []*Repository
	resp, err := s.client.Do(req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}

func (s *RepositoriesService) ListRepositoryUnits(repository string) ([]*Unit, *Response, error) {
	// set fields (must also be defined in the unit struct)
	fields := []string{
		"name",
		"version",
		"filename",
		"requires",
	}

	s.client.Units.SetFields(fields)
	r, resp, err := s.client.Units.ListUnits(repository)
	return r, resp, err
}

type GetRepositoryOptions struct {
	Details bool `url:"details,omitempty" json:"details,omitempty"`
}

func (s *RepositoriesService) GetRepository(
	repository string,
	opt *GetRepositoryOptions) (*Repository, *Response, error) {
	u := fmt.Sprintf("repositories/%s/", repository)

	req, err := s.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	r := new(Repository)
	resp, err := s.client.Do(req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}

func (s *RepositoriesService) SyncRepository(repository string) (*CallReport, *Response, error) {
	u := fmt.Sprintf("repositories/%s/actions/sync/", repository)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	cr := new(CallReport)
	resp, err := s.client.Do(req, &cr)

	if err != nil {
		return nil, resp, err
	}

	return cr, resp, err
}

type CopyRepositoryUnitsOptions struct {
	SourceRepository string                   `json:"source_repo_id,omitempty"`
	Criteria         *UnitAssociationCriteria `json:"criteria,omitempty"`
}

func (s *RepositoriesService) CopyRepositoryUnits(
	source_repository_id string,
	destination_repo_id string,
	criteria *UnitAssociationCriteria,
	// queryMap map[string][]map[string]map[string]string,
) (
	*CallReport,
	*Response,
	error,
) {

	u := fmt.Sprintf("repositories/%s/actions/associate/", destination_repo_id)

	opt := &CopyRepositoryUnitsOptions{
		SourceRepository: source_repository_id,
		Criteria:         criteria,
	}

	jsonOpt, err := json.Marshal(opt)
	fmt.Printf("JsonOpt: %s\n", jsonOpt)

	req, err := s.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	cr := new(CallReport)
	resp, err := s.client.Do(req, &cr)
	if err != nil {
		return nil, resp, err
	}

	return cr, resp, err
}
