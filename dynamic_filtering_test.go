package dqk

import (
	"encoding/json"
	"fmt"
	"maps"
	"strings"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

func TestIsAggregate(t *testing.T) {

	tests := []struct {
		name     string
		input    []Filters
		expected bool
	}{
		{
			name: "not a single aggregate",
			input: []Filters{
				{Name: "min dot id", Operator: "=", DbField: "min.id"},
				{Name: "max dot it", Operator: ">", DbField: "max.id"},
				{Name: "sum dot id", Operator: "<", DbField: "SUM.id"},
				{Name: "count in plain text", Operator: "<", DbField: "sum_prices"},
				{Name: "sum in plain text", Operator: "<", DbField: "sum_prices"},
				{Name: "avg_total", Operator: "<", DbField: "avg_total"},
				{Name: "total_avg", Operator: "<", DbField: "total_avg"},
				{Name: "min_sum", Operator: "<", DbField: "min_sum"},
				{Name: "MAX concatenated", Operator: "<", DbField: "MAXvalue"},
			},
			expected: false,
		},
		{
			name: "Aggregates",
			input: []Filters{
				{Name: "sum uppercase", Operator: "=", DbField: "SUM(price)"},
				{Name: "count uppercase", Operator: "=", DbField: "COUNT(price)"},
				{Name: "count multi case", Operator: "=", DbField: "CoUnT(price)"},
				{Name: "count lowercase", Operator: "=", DbField: "count(price)"},
				{Name: "min lowercase", Operator: "=", DbField: "min(price)"},
				{Name: "MAX uppercase", Operator: "=", DbField: "max(price)"},
				{Name: "max lowercase", Operator: "=", DbField: "MAX(price)"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for index, filter := range tt.input {
				assert.Equalf(t, tt.expected, filter.IsAggregate(), "Input: %s, index: %d", filter.DbField, index)

			}
		})
	}
}

func TestAreFiltersAggregate(t *testing.T) {
	tests := []struct {
		name     string
		input    []Filters
		expected bool
	}{
		{
			name: "not a single aggregate",
			input: []Filters{
				{Name: "a", Operator: "=", DbField: "min.id", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "max.id", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "SUM.id", FieldID: "3"},
			},
			expected: false,
		},
		{
			name: "with aggregate",
			input: []Filters{
				{Name: "a", Operator: "=", DbField: "min(id)", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "max(id)", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "sum(id)", FieldID: "3"},
			},
			expected: true,
		},
		{
			name: "with sum aggregate",
			input: []Filters{
				{Name: "a", Operator: "=", DbField: "min.id", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "max.id", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "sum(id)", FieldID: "3"},
			},
			expected: true,
		},
		{
			name: "with min aggregate",
			input: []Filters{
				{Name: "a", Operator: "=", DbField: "min(id)", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "max.id", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "sum.id", FieldID: "3"},
			},
			expected: true,
		},
		{
			name: "with max aggregate",
			input: []Filters{
				{Name: "a", Operator: "=", DbField: "min.id", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "max(id)", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "sum.id", FieldID: "3"},
			},
			expected: true,
		},
		{
			name: "with count aggregate",
			input: []Filters{
				{Name: "a", Operator: "=", DbField: "min.id", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "count(id)", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "sum.id", FieldID: "3"},
			},
			expected: true,
		},
		{
			name: "with stddev aggregate",
			input: []Filters{
				{Name: "a", Operator: "=", DbField: "min.id", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "stddev(id)", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "sum.id", FieldID: "3"},
			},
			expected: true,
		},
		{
			name: "with variance aggregate",
			input: []Filters{
				{Name: "a", Operator: "=", DbField: "min.id", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "variance(id)", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "sum.id", FieldID: "3"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AreFiltersAggregate(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}

}

func TestBuildFilterConditions(t *testing.T) {
	tests := []struct {
		name         string
		filters      []Filters
		values       map[string][]string
		Conditionals []Conditional
	}{
		{
			name: "present in filters only where claus",
			filters: []Filters{
				{Name: "a", Operator: "=", DbField: "min.id", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "max.id", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "SUM.id", FieldID: "3"},
			},
			values: map[string][]string{
				"a": {"5"},
				"b": {"5"},
				"c": {"8"},
			},
			Conditionals: []Conditional{
				NewConditional(
					sq.Expr("min.id = ?", "5"),
					TokenWhere,
					[]string{"5"},
				),
				NewConditional(
					sq.Expr("max.id > ?", "5"),
					TokenWhere,
					[]string{"5"},
				),
				NewConditional(
					sq.Expr("SUM.id < ?", "8"),
					TokenWhere,
					[]string{"8"},
				),
			}},
		{
			name: "present in filters only having claus",
			filters: []Filters{
				{Name: "a", Operator: "=", DbField: "MIN(id)", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "max(id)", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "SUM(id)", FieldID: "3"},
			},
			values: map[string][]string{
				"a": {"5"},
				"b": {"5"},
				"c": {"8"},
			},
			Conditionals: []Conditional{
				NewConditional(
					sq.Expr("MIN(id) = ?", "5"),
					TokenHaving,
					[]string{"5"},
				),
				NewConditional(
					sq.Expr("max(id) > ?", "5"),
					TokenHaving,
					[]string{"5"},
				),
				NewConditional(
					sq.Expr("SUM(id) < ?", "8"),
					TokenHaving,
					[]string{"8"},
				),
			},
		},
		{
			name: "present in filters both clauses",
			filters: []Filters{
				{Name: "a", Operator: "=", DbField: "MIN(id)", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "max(id)", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "SUM(id)", FieldID: "3"},
				{Name: "d", Operator: "=", DbField: "min.id", FieldID: "1"},
				{Name: "e", Operator: ">", DbField: "max.id", FieldID: "2"},
				{Name: "f", Operator: "<", DbField: "SUM.id", FieldID: "3"},
			},
			values: map[string][]string{
				"a":      {"5"},
				"b":      {"5"},
				"c":      {"8"},
				"d":      {"9"},
				"e":      {"1"},
				"f":      {"hello"},
				"limit":  {"100"},
				"offset": {"200"},
			},
			Conditionals: []Conditional{
				NewConditional(
					sq.Expr("MIN(id) = ?", "5"),
					TokenHaving,
					[]string{"5"},
				),
				NewConditional(
					sq.Expr("max(id) > ?", "5"),
					TokenHaving,
					[]string{"5"},
				),
				NewConditional(
					sq.Expr("SUM(id) < ?", "8"),
					TokenHaving,
					[]string{"8"},
				),
				NewConditional(
					sq.Expr("min.id = ?", "9"),
					TokenWhere,
					[]string{"9"},
				),
				NewConditional(
					sq.Expr("max.id > ?", "1"),
					TokenWhere,
					[]string{"1"},
				),
				NewConditional(
					sq.Expr("SUM.id < ?", "hello"),
					TokenWhere,
					[]string{"hello"},
				),
				NewConditional(
					sq.Expr("limit ?", "100"),
					TokenLimit,
					[]string{"100"},
				),
				NewConditional(
					sq.Expr("offset ?", "200"),
					TokenOffset,
					[]string{"200"},
				),
			},
		},
		{
			name: "LIKE filtering",
			filters: []Filters{
				{Name: "country", Operator: "LIKE", DbField: "country.name", FieldID: "1"},
			},
			values: map[string][]string{
				"country": {"Greece"},
			},
			Conditionals: []Conditional{
				NewConditional(
					sq.Expr("country.name LIKE ?", "%Greece%"),
					TokenWhere,
					[]string{"%Greece%"},
				),
			},
		},
		{
			name: "ILIKE filtering",
			filters: []Filters{
				{Name: "country", Operator: "ILIKE", DbField: "country.name", FieldID: "1"},
			},
			values: map[string][]string{
				"country": {"Greece"},
			},
			Conditionals: []Conditional{
				NewConditional(
					sq.Expr("country.name ILIKE ?", "%Greece%"),
					TokenWhere,
					[]string{"%Greece%"},
				),
			},
		},
		{
			name: "IN filtering",
			filters: []Filters{
				{Name: "country", Operator: "IN", DbField: "country.name", FieldID: "1"},
			},
			values: map[string][]string{
				"country": {"Greece", "France", "Germany"},
			},
			Conditionals: []Conditional{
				NewConditional(
					sq.Eq{"country.name": []string{"Greece", "France", "Germany"}},
					TokenWhere,
					[]string{"Greece", "France", "Germany"},
				),
			},
		},
		{
			name: "NULL & NOT NULL",
			filters: []Filters{
				{Name: "created_date", Operator: "=", DbField: "country.created_date", FieldID: "1"},
				{Name: "deleted_date", Operator: "=", DbField: "country.deleted_date", FieldID: "1"},
			},
			values: map[string][]string{
				"created_date": {"__NOT_NULL__", "France", "Germany"},
				"deleted_date": {"__NULL__", "France", "Germany"},
			},
			Conditionals: []Conditional{
				NewConditional(
					sq.Expr("country.created_date IS NOT NULL"),
					TokenWhere,
					[]string{"__NOT_NULL__", "France", "Germany"},
				),
				NewConditional(
					sq.Expr("country.deleted_date IS NULL"),
					TokenWhere,
					[]string{"__NULL__", "France", "Germany"},
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditionals := BuildFilterConditions(tt.filters, tt.values)
			assert.ElementsMatch(t, tt.Conditionals, conditionals)
		})

	}

}

func TestValidateParams(t *testing.T) {
	tests := []struct {
		name     string
		filters  []Filters
		values   map[string][]string
		Expected map[Filters][]string
	}{
		{
			name: "complete case testing",
			filters: []Filters{
				{Name: "country", Operator: "=", DbField: "country.name", FieldID: "1"},
				{Name: "stars", Operator: "IN", DbField: "booking.stars", FieldID: "1"},
				{Name: "deleted_date", Operator: "=", DbField: "country.deleted_date", FieldID: "1"},
				{Name: "created_date", Operator: "=", DbField: "country.created_date", FieldID: "1"},
				{Name: "c", Operator: "<", DbField: "SUM(id)", FieldID: "3"},
				{Name: "flying", Operator: "LIKE", DbField: "cars.flying", FieldID: "1"},
				{Name: "crying", Operator: "ILIKE", DbField: "cars.crying", FieldID: "1"},
			},
			values: map[string][]string{
				"country":      {"Greece"},
				"stars":        {"1", "2"},
				"deleted_date": {"__NULL__", "France", "Germany"},
				"created_date": {"__NOT_NULL__", "France", "Germany"},
				"c":            {"8"},
				"flying":       {"cars"},
				"crying":       {"TeSlA"},
				"limit":        {"100"},
				"offset":       {"200"},
			},
			Expected: map[Filters][]string{
				{Name: "country", Operator: "=", DbField: "country.name", FieldID: "1"}:              {"Greece"},
				{Name: "stars", Operator: "IN", DbField: "booking.stars", FieldID: "1"}:              {"1", "2"},
				{Name: "deleted_date", Operator: "=", DbField: "country.deleted_date", FieldID: "1"}: {"Not NULL"},
				{Name: "created_date", Operator: "=", DbField: "country.created_date", FieldID: "1"}: {"NULL"},
				{Name: "c", Operator: "<", DbField: "SUM(id)", FieldID: "3"}:                         {"8"},
				{Name: "flying", Operator: "LIKE", DbField: "cars.flying", FieldID: "1"}:             {"cars"},
				{Name: "crying", Operator: "ILIKE", DbField: "cars.crying", FieldID: "1"}:            {"TeSla"},
				{Name: "limit", Operator: "", DbField: "", FieldID: ""}:                              {"100"},
				{Name: "offset", Operator: "", DbField: "", FieldID: ""}:                             {"200"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := ValidateParams(tt.filters, tt.values)

			pa, _ := json.Marshal(params)
			pe, _ := json.Marshal(tt.Expected)

			fmt.Printf("expected: (%s) actual: (%s)", string(pe), string(pa))
			assert.Equal(t, string(pe), string(pa))

		})
	}
}

func TestDynamicFilters(t *testing.T) {
	tests := []struct {
		name                 string
		filters              []Filters
		values               map[string][]string
		ExpectedQuery        sq.SelectBuilder
		ExpectedStringParams map[string]string
	}{
		{
			name: "complete case testing",
			filters: []Filters{
				{Name: "country", Operator: "=", DbField: "country.name", FieldID: "1"},
				{Name: "stars", Operator: "IN", DbField: "booking.stars", FieldID: "1"},
				{Name: "deleted_date", Operator: "=", DbField: "country.deleted_date", FieldID: "1"},
				{Name: "created_date", Operator: "=", DbField: "country.created_date", FieldID: "1"},
				{Name: "c", Operator: "<", DbField: "SUM(id)", FieldID: "3"},
				{Name: "flying", Operator: "LIKE", DbField: "cars.flying", FieldID: "1"},
				{Name: "crying", Operator: "ILIKE", DbField: "cars.crying", FieldID: "1"},
			},
			values: map[string][]string{
				"country":      {"Greece"},
				"stars":        {"1", "2"},
				"deleted_date": {"__NULL__", "France", "Germany"},
				"created_date": {"__NOT_NULL__", "France", "Germany"},
				"c":            {"8"},
				"flying":       {"cars"},
				"crying":       {"TeSlA"},
			},
			ExpectedQuery: sq.Select("*").
				From("countries").
				Where(sq.Expr("country.deleted_date ?", "IS NULL")).
				Where(sq.Expr("country.created_date ?", "IS NOT NULL")).
				Where(sq.Like{"cars.flying": "cars"}).
				Where(sq.ILike{"cars.crying": "TeSlA"}).
				Where(sq.Eq{"country.name": "Greece"}).
				Where(sq.Eq{"booking.stars": "1,2"}).
				Having(sq.Expr("SUM(id) < ?", "8")),
			ExpectedStringParams: map[string]string{
				"country.name =":       "?",
				"booking.stars IN":     "(?,?)",
				"country.deleted_date": "IS NULL",
				"country.created_date": "IS NOT NULL",
				"SUM(id) <":            "?",
				"cars.flying LIKE":     "?",
				"cars.crying ILIKE":    "?",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			q := sq.Select("*").From("countries")
			v := maps.Clone(tt.values)

			query := DynamicFilters(tt.filters, q, v)

			ActualQueryStr, _, _ := query.ToSql()

			for name, value := range tt.ExpectedStringParams {
				fmt.Printf("Expected: (%s %s) Match: %t\n", name, value, strings.Contains(ActualQueryStr, fmt.Sprintf("%s %s", name, value)))
				assert.True(t, strings.Contains(ActualQueryStr, fmt.Sprintf("%s %s", name, value)))

			}
			fmt.Println(ActualQueryStr)

			assert.Equal(t, tt.values["country"][0], "Greece")
			assert.Equal(t, tt.values["stars"][0], "1")
			assert.Equal(t, tt.values["stars"][1], "2")
			assert.Equal(t, tt.values["c"][0], "8")
			assert.Equal(t, tt.values["flying"][0], "%cars%")
			assert.Equal(t, tt.values["crying"][0], "%TeSlA%")

		})

	}

}

func TestExtendFilters(t *testing.T) {
	tests := []struct {
		name     string
		f1       []Filters
		f2       []Filters
		expected []Filters
	}{
		{
			name: "f1 longer than f2",
			f1: []Filters{
				{Name: "a", Operator: "=", DbField: "field1", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "field2", FieldID: "2"},
			},
			f2: []Filters{
				{Name: "c", Operator: "<", DbField: "field3", FieldID: "3"},
			},
			expected: []Filters{
				{Name: "c", Operator: "<", DbField: "field3", FieldID: "3"},
				{Name: "a", Operator: "=", DbField: "field1", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "field2", FieldID: "2"},
			},
		},
		{
			name: "f2 longer than f1",
			f1: []Filters{
				{Name: "x", Operator: "!=", DbField: "field9", FieldID: "9"},
			},
			f2: []Filters{
				{Name: "y", Operator: "LIKE", DbField: "field8", FieldID: "8"},
				{Name: "z", Operator: "<=", DbField: "field7", FieldID: "7"},
			},
			expected: []Filters{
				{Name: "y", Operator: "LIKE", DbField: "field8", FieldID: "8"},
				{Name: "z", Operator: "<=", DbField: "field7", FieldID: "7"},
				{Name: "x", Operator: "!=", DbField: "field9", FieldID: "9"},
			},
		},
		{
			name:     "both empty",
			f1:       []Filters{},
			f2:       []Filters{},
			expected: []Filters{},
		},
		{
			name: "equal length",
			f1: []Filters{
				{Name: "e1", Operator: "=", DbField: "f1", FieldID: "id1"},
			},
			f2: []Filters{
				{Name: "e2", Operator: "!=", DbField: "f2", FieldID: "id2"},
			},
			expected: []Filters{
				{Name: "e2", Operator: "!=", DbField: "f2", FieldID: "id2"},
				{Name: "e1", Operator: "=", DbField: "f1", FieldID: "id1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := [][]Filters{
				tt.f2,
				tt.f1,
			}
			result := ExtendFilters(f)
			assert.Equal(t, tt.expected, result)
		})
	}
}
