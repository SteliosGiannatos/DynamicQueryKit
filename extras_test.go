package dynamicquerykit

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

func TestGetPaginationQuery(t *testing.T) {
	tests := []struct {
		name          string
		query         sq.SelectBuilder
		expectedQuery string
		expectedArgs  []any
	}{
		{
			name:          "with args",
			query:         sq.Select("id").From("cars").Where(sq.Eq{"cars.color": "black"}),
			expectedQuery: "SELECT COUNT(*) AS total_rows FROM ( SELECT id FROM cars WHERE cars.color = ? ) AS grouped_results;",
			expectedArgs:  []any{"black"},
		},
		{
			name:          "no args",
			query:         sq.Select("id").From("cars"),
			expectedQuery: "SELECT COUNT(*) AS total_rows FROM ( SELECT id FROM cars ) AS grouped_results;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := GetPaginationQuery(tt.query)
			query, _, _ := val.ToSql()
			assert.Equal(t, tt.expectedQuery, query)
		})

	}

}
