package dynamicquerykit

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

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
				{Name: "a", Operator: "=", DbField: "field1", FieldID: "1"},
				{Name: "b", Operator: ">", DbField: "field2", FieldID: "2"},
				{Name: "c", Operator: "<", DbField: "field3", FieldID: "3"},
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
			result := ExtendFilters(tt.f1, tt.f2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAggregate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "SUM with parentheses", input: "SUM(price)", expected: true},
		{name: "COUNT with parentheses", input: "COUNT(price)", expected: true},
		{name: "count in plain text", input: "count_price", expected: false},
		{name: "sum in plain text", input: "sum_prices", expected: false},
		{name: "avg_total", input: "avg_total", expected: false},
		{name: "total_avg", input: "total_avg", expected: false},
		{name: "min_sum", input: "min_sum", expected: false},
		{name: "min parentheses", input: "min(total_price)", expected: true},
		{name: "MIN uppercase", input: "MIN(value)", expected: true},
		{name: "max parentheses", input: "max(total_price)", expected: true},
		{name: "MAX uppercase", input: "MAX(value)", expected: true},
		{name: "MAX concatenated", input: "MAXvalue", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAggregate(tt.input)
			// comment at the fix:
			// changed assert.Equal(...) to assert.Equalf(...) to include formatted message
			assert.Equalf(t, tt.expected, result, "Input: %s", tt.input) // <-- fix
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
		name           string
		filters        []Filters
		values         map[string][]string
		expectedWhere  []sq.Sqlizer
		expectedHaving []sq.Sqlizer
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
			expectedWhere: []sq.Sqlizer{
				sq.Expr("min.id = ?", "5"),
				sq.Expr("max.id > ?", "5"),
				sq.Expr("SUM.id < ?", "8"),
			},
			expectedHaving: []sq.Sqlizer{},
		},
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
			expectedWhere: []sq.Sqlizer{},
			expectedHaving: []sq.Sqlizer{
				sq.Expr("MIN(id) = ?", "5"),
				sq.Expr("max(id) > ?", "5"),
				sq.Expr("SUM(id) < ?", "8"),
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
				"a": {"5"},
				"b": {"5"},
				"c": {"8"},
				"d": {"9"},
				"e": {"1"},
				"f": {"hello"},
			},
			expectedWhere: []sq.Sqlizer{
				sq.Expr("min.id = ?", "9"),
				sq.Expr("max.id > ?", "1"),
				sq.Expr("SUM.id < ?", "hello"),
			},
			expectedHaving: []sq.Sqlizer{
				sq.Expr("MIN(id) = ?", "5"),
				sq.Expr("max(id) > ?", "5"),
				sq.Expr("SUM(id) < ?", "8"),
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
			expectedWhere: []sq.Sqlizer{
				sq.Expr("country.name LIKE ?", "%Greece%"),
			},
			expectedHaving: []sq.Sqlizer{},
		},
		{
			name: "ILIKE filtering",
			filters: []Filters{
				{Name: "country", Operator: "ILIKE", DbField: "country.name", FieldID: "1"},
			},
			values: map[string][]string{
				"country": {"Greece"},
			},
			expectedWhere: []sq.Sqlizer{
				sq.Expr("country.name ILIKE ?", "%Greece%"),
			},
			expectedHaving: []sq.Sqlizer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResultsWhere, ResultsHaving := BuildFilterConditions(tt.filters, tt.values)
			assert.Equal(t, tt.expectedWhere, ResultsWhere)
			assert.Equal(t, tt.expectedHaving, ResultsHaving)
		})

	}

}
