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

func (s *RepositoriesService) ListRepositories() ([]*Repository, *Response, error) {
	req, err := s.client.NewRequest("GET", "repositories/", nil)
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

type CreateRepositoryOptions struct {
	Details bool `url:"details,omitempty" json:"details,omitempty"`
}

func (s *RepositoriesService) GetRepository(
	repository string,
	opt *CreateRepositoryOptions) (*Repository, *Response, error) {
	u := fmt.Sprintf("repositories/%s/", repository)

	req, err := s.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	r := new(Repository)
	resp, err := s.client.Do(req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}

type Importer struct {
	Id             string          `json:"id"`
	ImporterConfig *ImporterConfig `json:"config"`
}

type ImporterConfig struct {
	Feed          string `json:"feed"`
	RemoveMissing bool   `json:"remove_missing"`
}
