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

type UnitsService struct {
	client *Client
	Fields []string
}

func (s *UnitsService) SetFields(fields []string) {
	s.Fields = fields
}

type Unit struct {
	Id       string `json:"id"`
	RepoId   string `json:"repo_id"`
	TypeId   string `json:"unit_type_id"`
	UnitId   string `json:"unit_id"`
	Metadata struct {
		Name     string    `json:"name"`
		Version  string    `json:"version"`
		FileName string    `json:"filename"`
		Requires []Require `json:"requires"`
	} `json:"metadata"`
}

type Require struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

//  Options
type GetUnitOptions struct {
	*SearchCriteria `json:"criteria,omitempty"`
}

func (s *UnitsService) GetUnitByTaskResult(result *Result) ([]*Unit, *Response, error) {

	orFilter := NewFilter("$or")
	for _, resultUnit := range result.ResultUnits {
		andFilter := NewFilter("$and")

		andFilter.AddExpression("name", "$regex", fmt.Sprintf("^%s$", resultUnit.UnitKey.Name))
		andFilter.AddExpression("version", "$regex", fmt.Sprintf("^%s$", resultUnit.UnitKey.Version))
		andFilter.AddExpression("release", "$regex", fmt.Sprintf("^%s$", resultUnit.UnitKey.Release))
		andFilter.AddExpression("arch", "$regex", fmt.Sprintf("^%s$", resultUnit.UnitKey.Arch))

		orFilter.AddSubFilter(andFilter)
	}

	criteria := NewSearchCriteria()
	criteria.AddField("filename")
	criteria.AddFilter(orFilter)

	opt := GetUnitOptions{
		SearchCriteria: criteria,
	}

	// url := fmt.Sprintf("content/units/%s/search/", resultUnit.TypeId)
	url := fmt.Sprintf("content/units/%s/search/", "rpm")
	req, err := s.client.NewRequest("POST", url, opt)
	if err != nil {
		return nil, nil, err
	}

	var u []*Unit
	resp, err := s.client.Do(req, &u)
	if err != nil {
		return nil, resp, err
	}

	return u, resp, err
}

//  Options
type ListUnitsOptions struct {
	*UnitAssociationCriteria `json:"criteria,omitempty"`
}

func (s *UnitsService) ListUnits(repository string) ([]*Unit, *Response, error) {
	// units options

	criteria := NewUnitAssociationCriteria()
	criteria.AddFields(s.Fields)

	opt := ListUnitsOptions{
		UnitAssociationCriteria: criteria,
	}

	url := fmt.Sprintf("repositories/%s/search/units/", repository)
	req, err := s.client.NewRequest("POST", url, opt)
	if err != nil {
		return nil, nil, err
	}

	var u []*Unit
	resp, err := s.client.Do(req, &u)
	if err != nil {
		return nil, resp, err
	}

	return u, resp, err
}
