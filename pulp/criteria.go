package pulp

import (
	"github.com/linkosmos/mapop"
)

// Criteria

type UnitAssociationCriteria struct {
	UnitAssociationFilters `json:"filters,omitempty"`
	UnitAssociationFields  `json:"fields,omitempty"`
}

func NewUnitAssociationCriteria() (criteria *UnitAssociationCriteria) {
	criteria = &UnitAssociationCriteria{}
	return
}

func (uac *UnitAssociationCriteria) AddFilter(filter *Filter) {
	uac.UnitAssociationFilters.AddFilter(filter)
}

func (uac *UnitAssociationCriteria) AddFilterMap(filterMap map[string]interface{}) {
	uac.UnitAssociationFilters.AddFilterMap(filterMap)
}

func (uac *UnitAssociationCriteria) AddFields(fields []string) {
	uac.UnitAssociationFields.AddFields(fields)
}

func (uac *UnitAssociationCriteria) AddField(field string) {
	uac.UnitAssociationFields.AddField(field)
}

type SearchCriteria struct {
	SearchFilters map[string]interface{} `json:"filters,omitempty"`
	SearchFields  []string               `json:"fields,omitempty"`
}

func NewSearchCriteria() (criteria *SearchCriteria) {
	criteria = &SearchCriteria{}
	return
}

func (sc *SearchCriteria) AddFilter(filter *Filter) {
	sc.AddFilterMap(filter.GetMap())
}

func (sc *SearchCriteria) AddFilterMap(filterMap map[string]interface{}) {
	sc.SearchFilters = mapop.Merge(sc.SearchFilters, filterMap)
}

func (sc *SearchCriteria) AddFields(fields []string) {
	for _, field := range fields {
		sc.AddField(field)
	}
}

func (sc *SearchCriteria) AddField(field string) {
	sc.SearchFields = append(sc.SearchFields, field)
}

// Unit association fields

type UnitAssociationFields struct {
	Unit        []string `json:"unit,omitempty"`
	Association []string `json:"association,omitempty"`
}

func NewUnitAssociationFields(fields []string) (unitAssociationFields UnitAssociationFields) {
	unitAssociationFields = UnitAssociationFields{
		Unit: fields,
	}
	return
}

// Unit association filters

type UnitAssociationFilters struct {
	Unit        map[string]interface{} `json:"unit,omitempty"`
	Association []string               `json:"association,omitempty"`
}

func (uaf *UnitAssociationFilters) AddFilter(filter *Filter) {
	uaf.AddFilterMap(filter.GetMap())
}

func (uaf *UnitAssociationFilters) AddFilterMap(filterMap map[string]interface{}) {
	uaf.Unit = mapop.Merge(uaf.Unit, filterMap)
}

func (uafields *UnitAssociationFields) AddFields(fields []string) {
	for _, field := range fields {
		uafields.AddField(field)
	}
}

func (uafields *UnitAssociationFields) AddField(field string) {
	uafields.Unit = append(uafields.Unit, field)
}

func NewUnitAssociationFiters(filters map[string]interface{}) (unitAssociationFilters UnitAssociationFilters) {
	unitAssociationFilters = UnitAssociationFilters{
		Unit: filters,
	}
	return
}

// Common filter

type Filter struct {
	Operator    string
	Expressions []*Expression
	SubFilters  []*Filter
}

func (f *Filter) AddExpression(field string, selector string, value string) {
	expression := &Expression{
		UnitField: field,
		Selector:  selector,
		Value:     value,
	}
	f.Expressions = append(f.Expressions, expression)
}

func NewFilter(operator string) (filter *Filter) {
	filter = &Filter{
		Operator: operator,
	}
	return
}

func (f *Filter) AddSubFilter(filter *Filter) {
	f.SubFilters = append(f.SubFilters, filter)
}

func (f *Filter) GetMap() (filterMap map[string]interface{}) {
	expressions := []map[string]map[string]string{}
	for _, exp := range f.Expressions {
		expressions = append(expressions, exp.GetMap())
	}
	filterMap = map[string]interface{}{
		f.Operator: expressions,
	}

	return
}

// Expression

type Expression struct {
	UnitField string
	Selector  string
	Value     string
}

func NewExpression(field string, selector string, value string) (filter *Expression) {
	filter = &Expression{
		UnitField: field,
		Selector:  selector,
		Value:     value,
	}
	return
}

func (f *Expression) GetMap() (filterMap map[string]map[string]string) {
	filterMap = map[string]map[string]string{
		f.UnitField: {
			f.Selector: f.Value,
		},
	}
	return
}

func MergeFilterMap(
	filterMap1 map[string]interface{},
	filterMap2 map[string]interface{},
) (mergedFilterMap map[string]interface{}) {
	mergedFilterMap = mapop.Merge(filterMap1, filterMap2)
	return
}
