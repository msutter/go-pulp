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

type UnitsService struct {
	client *Client
	Fields []string
}

func (s *UnitsService) SetFields(fields []string) {
	s.Fields = fields
}

type Unit struct {
	UnitId   string    `json:"_id"`
	TypeId   string    `json:"_content_type_id"`
	FileName string    `json:"filename"`
	Name     string    `json:"name"`
	Version  string    `json:"version"`
	Release  string    `json:"release"`
	Arch     string    `json:"arch"`
	Requires []Require `json:"requires"`
}

type SearchUnit struct {
	UnitId                string     `json:"_id"`
	Name                  string     `json:"name"`
	Version               string     `json:"version"`
	Release               string     `json:"release"`
	Arch                  string     `json:"arch"`
	Epoch                 string     `json:"epoch"`
	FileName              string     `json:"filename"`
	Requires              []*Require `json:"requires"`
	RepositoryMemberships []string   `json:"repository_memberships"`
}

type ContentUnitAssociation struct {
	Id       string `json:"id"`
	RepoId   string `json:"repo_id"`
	TypeId   string `json:"unit_type_id"`
	UnitId   string `json:"unit_id"`
	Metadata struct {
		Name     string     `json:"name"`
		Version  string     `json:"version"`
		FileName string     `json:"filename"`
		Requires []*Require `json:"requires"`
	} `json:"metadata"`
}

type Require struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Release string `json:"release"`
	Epoch   string `json:"epoch"`
	Flags   string `json:"flags"`
}

//  Options
type GetUnitOptions struct {
	*SearchCriteria `json:"criteria,omitempty"`
}

func (s *UnitsService) GetUnitFileNamesByTaskResult(result *Result) (filenames []string, err error) {
	var units []*Unit
	units, _, err = s.GetUnitsByTaskResult(result, []string{"filename"})
	for _, unit := range units {
		filenames = append(filenames, unit.FileName)
	}
	return filenames, err
}

func (s *UnitsService) GetUnitsByTaskResult(result *Result, fields []string) ([]*Unit, *Response, error) {

	var u []*Unit
	if len(result.ResultUnits) > 0 {
		resultType := result.ResultUnits[0].TypeId
		orFilter := NewFilter()
		orFilter.Operator = "$or"
		for _, resultUnit := range result.ResultUnits {

			// Unit filter (no need for filter operator)
			unitFilter := NewFilter()
			unitFilter.AddExpression("name", "$regex", fmt.Sprintf("^%s$", resultUnit.UnitKey.Name))
			unitFilter.AddExpression("version", "$regex", fmt.Sprintf("^%s$", resultUnit.UnitKey.Version))
			unitFilter.AddExpression("release", "$regex", fmt.Sprintf("^%s$", resultUnit.UnitKey.Release))
			unitFilter.AddExpression("arch", "$regex", fmt.Sprintf("^%s$", resultUnit.UnitKey.Arch))

			orFilter.AddSubFilter(unitFilter)
		}

		criteria := NewSearchCriteria()
		for _, field := range fields {
			criteria.AddField(field)
		}
		criteria.AddFilter(orFilter)

		opt := GetUnitOptions{
			SearchCriteria: criteria,
		}

		// // check request body
		// jsonSearchOpt, err := json.Marshal(opt)
		// fmt.Printf("jsonSearchOpt: %s\n", jsonSearchOpt)

		url := fmt.Sprintf("content/units/%s/search/", resultType)
		req, err := s.client.NewRequest("POST", url, opt)
		if err != nil {
			return nil, nil, err
		}

		resp, err := s.client.Do(req, &u)
		if err != nil {
			return nil, resp, err
		}

		return u, resp, err
	} else {
		return u, nil, nil
	}

}

//  Options
type ListUnitsOptions struct {
	*UnitAssociationCriteria `json:"criteria,omitempty"`
}

func (s *UnitsService) ListUnits(repository string) ([]*ContentUnitAssociation, *Response, error) {
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

	var u []*ContentUnitAssociation
	resp, err := s.client.Do(req, &u)
	if err != nil {
		return nil, resp, err
	}

	return u, resp, err
}

//  Options
type SearchUnitsOptions struct {
	*SearchCriteria `json:"criteria,omitempty"`
	IncludeRepos    bool `json:"include_repos,omitempty"`
}

func (s *UnitsService) SearchUnits(contentType string, criteria *SearchCriteria) ([]*SearchUnit, *Response, error) {

	opt := SearchUnitsOptions{
		SearchCriteria: criteria,
		IncludeRepos:   true,
	}

	url := fmt.Sprintf("content/units/%s/search/", contentType)
	req, err := s.client.NewRequest("POST", url, opt)
	if err != nil {
		return nil, nil, err
	}

	// check request body
	jsonSearchOpt, err := json.Marshal(opt)
	fmt.Printf("jsonSearchOpt: %s\n", jsonSearchOpt)

	var u []*SearchUnit

	fmt.Printf("\nDebug: Before request\n")
	resp, err := s.client.Do(req, &u)
	if err != nil {
		fmt.Printf("\nDebug: Error by request\n")
		fmt.Printf("\n%v\n",err.Error())
		return nil, resp, err
	}
	fmt.Printf("\nDebug: After request\n")

	return u, resp, err
}
