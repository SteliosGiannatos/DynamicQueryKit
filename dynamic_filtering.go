package dqk

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

// AreFiltersAggregate returns true if any of the provided filters use an allowedAggregateFunctions
func AreFiltersAggregate(filters []Filters) bool {
	aggregate := false
	for _, filter := range filters {
		if filter.IsAggregate() {
			aggregate = true
		}
	}
	return aggregate
}

func ValidateParams(filters []Filters, params map[string][]string) map[Filters][]string {
	conditionsSet := make(map[Filters][]string)

	for _, filter := range filters {
		values, ok := params[filter.Name]
		if len(values) == 0 || !ok {
			continue
		}
		filter.ApplyNullToken(values...)
		conditionsSet[filter] = values
	}

	return ValidateValus(conditionsSet)

}

func ValidateValus(filterValues map[Filters][]string) map[Filters][]string {
	for filter, values := range filterValues {
		for index, value := range values {
			if value == "" {
				continue
			}
			switch filter.Operator {
			case "LIKE", "ILIKE":
				values[index] = fmt.Sprintf("%%%s%%", value)
			}
		}
	}
	return filterValues
}

func GetParamsApplied(filters []Filters, params map[string][]string) map[string]string {

	filterValues := ValidateParams(filters, params)

	m := make(map[string]string)

	for filter, allowedValues := range filterValues {

		HasNullOrNotNull := filter.HasNullOrNotNull(allowedValues...)
		if HasNullOrNotNull {
			m[filter.DbField] = filter.Operator
			continue
		}

		if filter.Operator == "IN" && !HasNullOrNotNull {
			m[fmt.Sprintf("%s %s", filter.DbField, filter.Operator)] = strings.Join(allowedValues, ",")
			continue
		}
		for _, value := range allowedValues {
			m[fmt.Sprintf("%s %s", filter.DbField, filter.Operator)] = value
		}
	}

	return m
}

// BuildFilterConditions takes in allowed filters and values to be filtered. The key of the values map must match
// the Filter.Name field. It returns first all where conditions (conditions that should be added in a where claus)
// and having conditions (all conditions that should be added in a having claus) Both can be consolidated using
// either sq.And() or sq.Or() or a custom method in order to be applied to a filter
func BuildFilterConditions(filters []Filters, params map[string][]string) ([]sq.Sqlizer, []sq.Sqlizer) {
	filterValues := ValidateParams(filters, params)

	var (
		WhereConditions  = []sq.Sqlizer{}
		HavingConditions = []sq.Sqlizer{}
	)

	for filter, allowedValues := range filterValues {

		HasNullOrNotNull := filter.HasNullOrNotNull(allowedValues...)
		if HasNullOrNotNull {
			WhereConditions = append(WhereConditions, sq.Expr(fmt.Sprintf("%s %s", filter.DbField, filter.Operator)))
			continue
		}
		if filter.Operator == "IN" && !HasNullOrNotNull {
			WhereConditions = append(WhereConditions, sq.Eq{filter.DbField: allowedValues})
			continue
		}
		for _, value := range allowedValues {
			if filter.IsAggregate() {
				HavingConditions = append(HavingConditions, sq.Expr(fmt.Sprintf("%s %s ?", filter.DbField, filter.Operator), value))
				continue
			}
			WhereConditions = append(WhereConditions, sq.Expr(fmt.Sprintf("%s %s ?", filter.DbField, filter.Operator), value))
		}
	}
	return WhereConditions, HavingConditions
}

// DynamicFilters it applies dynamic filters based on the allowed filters. These are added to the specified query
// it can get the query params as is from the r.URL.query() method.
// it does not stop the user from passing multiple = params
// all conditions are passed as AND parameters. This is true for both having & where conditions
func DynamicFilters(f []Filters, q sq.SelectBuilder, queryParams map[string][]string) sq.SelectBuilder {
	whereCond, HavingCond := BuildFilterConditions(f, queryParams)
	if len(HavingCond) > 0 {
		q = q.Having(sq.And(HavingCond))
	}
	if len(whereCond) > 0 {
		q = q.Where(sq.And(whereCond))
	}
	return q
}

// ExtendFilters takes in n filters and returns a complete filter list
func ExtendFilters(filters [][]Filters) []Filters {
	combination := []Filters{}
	for _, filterList := range filters {
		combination = append(combination, filterList...)
	}

	return combination
}
