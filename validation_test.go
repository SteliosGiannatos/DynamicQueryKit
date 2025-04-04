package dynamicquerykit

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFieldJSONTag(t *testing.T) {
	type testStruct struct {
		ID   int    `json:"id" xml:"id" yaml:"id" csv:"id"`
		Name string `json:"name" xml:"name" yaml:"name" csv:"name"`
	}

	tests := []struct {
		name     string
		field    string
		expected bool
	}{
		{
			name:     "contained in struct",
			field:    "id",
			expected: true,
		},
		{
			name:     "not contained in struct",
			field:    "a",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsFieldJSONTag(testStruct{}, tt.field)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestIsFieldFilter(t *testing.T) {
	tests := []struct {
		name     string
		filters  []Filters
		field    string
		expected bool
	}{
		{
			name:  "contained in filters",
			field: "id",
			filters: []Filters{
				{Name: "name", Operator: "=", DbField: "table.name", FieldID: "table.name"},
				{Name: "id", Operator: "=", DbField: "table.id", FieldID: "table.id"},
			},
			expected: true,
		},
		{
			name:  "not contained in filters",
			field: "a",
			filters: []Filters{
				{Name: "name", Operator: "=", DbField: "table.name", FieldID: "table.name"},
				{Name: "id", Operator: "=", DbField: "table.id", FieldID: "table.id"},
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsFieldFilter(tt.filters, tt.field)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestOrderValidation(t *testing.T) {
	tests := []struct {
		name           string
		filters        []Filters
		orderBy        string
		OrderDirection string
		expected       string
	}{
		{
			name:           "in filter with provided direction",
			orderBy:        "id",
			OrderDirection: "desc",
			filters: []Filters{
				{Name: "id", Operator: "=", DbField: "table.id", FieldID: "table.id"},
			},
			expected: "table.id DESC",
		},
		{
			name:    "in filter without provided direction",
			orderBy: "id",
			filters: []Filters{
				{Name: "id", Operator: "=", DbField: "table.id", FieldID: "table.id"},
			},
			expected: "table.id ASC",
		},
		{
			name:           "not in filter with provided direction",
			orderBy:        "id",
			OrderDirection: "desc",
			filters: []Filters{
				{Name: "name", Operator: "=", DbField: "table.name", FieldID: "table.name"},
			},
			expected: "table.name DESC",
		},
		{
			name:    "not in filter without provided direction",
			orderBy: "id",
			filters: []Filters{
				{Name: "name", Operator: "=", DbField: "table.name", FieldID: "table.name"},
			},
			expected: "table.name ASC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OrderValidation(tt.orderBy, tt.OrderDirection, tt.filters)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestDatabaseValidation(t *testing.T) {
	originalLogger := slog.Default()
	f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("failed to open os.DevNull: %v", err)
	}
	defer f.Close()

	slog.SetDefault(slog.New(slog.NewTextHandler(f, nil)))

	defer slog.SetDefault(originalLogger)
	tests := []struct {
		name               string
		providedError      error
		expectedStatusCode int
	}{
		{
			name:               "conflict",
			providedError:      fmt.Errorf("1062"),
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:               "invalid related resource",
			providedError:      fmt.Errorf("1452"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "data too long for column",
			providedError:      fmt.Errorf("1406"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "null column provided",
			providedError:      fmt.Errorf("1048"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "no data returned",
			providedError:      fmt.Errorf("no rows in result set"),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "undefined error",
			providedError:      fmt.Errorf("some undefined error"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := DatabaseValidation(tt.providedError)
			assert.Equal(t, tt.expectedStatusCode, got)
		})
	}

}
