package dynamicquerykit

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

// AllowedAggregateFunctions contains all functions that can be used for filtering
var AllowedAggregateFunctions = []string{
	"count",
	"sum",
	"min",
	"max",
	"stddev",
	"variance",
}

// IsAggregate checks if a databases field is one of the allowed aggregate functions
func IsAggregate(field string) bool {
	lower := strings.ToLower(field)
	for _, agg := range AllowedAggregateFunctions {
		if strings.HasPrefix(lower, agg+"(") {
			return true
		}
	}
	return false
}

// AreFiltersAggregate returns true if any of the provided filters use an allowedAggregateFunctions
func AreFiltersAggregate(filters []Filters) bool {
	method := "AreFiltersAggregate"
	aggregate := false
	slog.LogAttrs(context.Background(), slog.LevelDebug, "checking aggregation for filters", slog.Any("filters", filters))
	for index, filter := range filters {
		if IsAggregate(filter.DbField) {
			aggregate = true
			slog.LogAttrs(context.Background(), slog.LevelDebug, "checking aggregation for filter", slog.String("method", method), slog.Int("aggregate filter index", index), slog.Any("filter", filter))
		}
	}
	return aggregate
}

// BuildFilterConditions takes in allowed filters and values to be filtered. The key of the values map must match
// the Filter.Name field. It returns first all where conditions (conditions that should be added in a where claus)
// and having conditions (all conditions that should be added in a having claus) Both can be consolidated using
// either sq.And() or sq.Or() or a custom method in order to be applied to a filter
func BuildFilterConditions(filters []Filters, values map[string][]string) ([]sq.Sqlizer, []sq.Sqlizer) {
	var (
		WhereConditions  = []sq.Sqlizer{}
		HavingConditions = []sq.Sqlizer{}
	)
	for _, filter := range filters {
		allowedValues := values[filter.Name]
		if len(allowedValues) == 0 {
			continue
		}
		useHaving := IsAggregate(filter.DbField)
		hasNullOrNot := false
		for index, value := range allowedValues {
			if value == "" {
				continue
			}
			switch value {
			case "__NULL__":
				filter.Operator = "IS NULL"
				hasNullOrNot = true
				continue
			case "__NOT_NULL__":
				filter.Operator = "IS NOT NULL"
				hasNullOrNot = true
				continue
			}
			switch filter.Operator {
			case "LIKE", "ILIKE":
				allowedValues[index] = fmt.Sprintf("%%%s%%", value)
			case "IN":
				continue
			}
		}
		if hasNullOrNot {
			WhereConditions = append(WhereConditions, sq.Expr(fmt.Sprintf("%s %s", filter.DbField, filter.Operator)))
			continue
		}
		if filter.Operator == "IN" && !hasNullOrNot {
			WhereConditions = append(WhereConditions, sq.Eq{filter.DbField: allowedValues})
			continue
		}
		for _, value := range allowedValues {
			if useHaving {
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

// ExtendFilters takes in two filters and appends one to the other
// it always appends the smallest filter to the larger filter.
// if both are equal size f1 gets appended to f2
func ExtendFilters(f1 []Filters, f2 []Filters) []Filters {
	if len(f1) > len(f2) {
		f1 = append(f1, f2...)
		return f1
	}
	f2 = append(f2, f1...)
	return f2
}
