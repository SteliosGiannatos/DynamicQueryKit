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
	newMap := map[string][]string{}
	for k, v := range params {
		newMap[strings.ToLower(k)] = v
	}

	limitOffset := []Filters{
		{Name: TokenLimit, Operator: "", DbField: "", FieldID: ""},
		{Name: TokenOffset, Operator: "", DbField: "", FieldID: ""},
	}
	filters = append(filters, limitOffset...)
	conditionsSet := make(map[Filters][]string)

	for _, filter := range filters {
		values, ok := newMap[strings.ToLower(filter.Name)]
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
		if len(allowedValues) <= 0 {
			continue
		}

		HasNullOrNotNull := filter.HasNullOrNotNull(allowedValues...)
		if HasNullOrNotNull {
			m[filter.DbField] = filter.Operator
			continue
		}

		if filter.Operator == "IN" && !HasNullOrNotNull {
			m[fmt.Sprintf("%s %s", filter.DbField, filter.Operator)] = strings.Join(allowedValues, ",")
			continue
		}

		if filter.Name == TokenLimit || filter.Name == TokenOffset {
			m[fmt.Sprintf("%s", filter.Name)] = allowedValues[0]
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
func BuildFilterConditions(filters []Filters, params map[string][]string) []Conditional {
	filterValues := ValidateParams(filters, params)
	var conditionals []Conditional

	for filter, values := range filterValues {
		if len(values) <= 0 {
			continue
		}
		HasNullOrNotNull := filter.HasNullOrNotNull(values...)
		if HasNullOrNotNull {
			conditionals = append(conditionals, NewConditional(
				sq.Expr(fmt.Sprintf("%s %s", filter.DbField, filter.Operator)),
				TokenWhere,
				values,
			))
			continue
		}

		if filter.Operator == "IN" && !HasNullOrNotNull {
			conditionals = append(conditionals, NewConditional(
				sq.Eq{filter.DbField: values},
				TokenWhere,
				values,
			))
			continue
		}

		if filter.Name == TokenLimit || filter.Name == TokenOffset {
			conditionals = append(conditionals, NewConditional(
				sq.Expr(fmt.Sprintf("%s ?", filter.Name), values[0]),
				filter.Name,
				values,
			))
			continue
		}

		for _, value := range values {
			if filter.IsAggregate() {
				conditionals = append(conditionals, NewConditional(
					sq.Expr(fmt.Sprintf("%s %s ?", filter.DbField, filter.Operator), value),
					TokenHaving,
					values,
				))
				continue
			}
			conditionals = append(conditionals, NewConditional(
				sq.Expr(fmt.Sprintf("%s %s ?", filter.DbField, filter.Operator), value),
				TokenWhere,
				values,
			))
		}
	}
	return conditionals
}

// DynamicFilters it applies dynamic filters based on the allowed filters. These are added to the specified query
// it can get the query params as is from the r.URL.query() method.
// it does not stop the user from passing multiple = params
// all conditions are passed as AND parameters. This is true for both having & where conditions
func DynamicFilters(f []Filters, q sq.SelectBuilder, queryParams map[string][]string) sq.SelectBuilder {
	conditions := BuildFilterConditions(f, queryParams)
	for _, condition := range conditions {
		q = condition.Apply(q)
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
