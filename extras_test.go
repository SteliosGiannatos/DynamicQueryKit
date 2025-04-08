package dqk

import (
	"fmt"
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

func TestGetRouteKey(t *testing.T) {
	id := 1
	tests := []struct {
		name             string
		route            string
		id               *int
		Args             []string
		expectedRouteKey string
		expectedIndexKey string
	}{
		{
			name:             "simple route",
			route:            "cars",
			id:               nil,
			Args:             []string{},
			expectedRouteKey: "cars",
			expectedIndexKey: "cars",
		},
		{
			name:             "route for asset",
			route:            "cars",
			id:               &id,
			Args:             []string{},
			expectedRouteKey: fmt.Sprintf("cars:%v", id),
			expectedIndexKey: fmt.Sprintf("cars:%v", id),
		},
		{
			name:             "route for asset and arg",
			route:            "cars",
			id:               &id,
			Args:             []string{"colors"},
			expectedRouteKey: fmt.Sprintf("cars:%v:colors", id),
			expectedIndexKey: fmt.Sprintf("cars:%v", id),
		},
		{
			name:             "route for asset and many args",
			route:            "cars",
			id:               &id,
			Args:             []string{"colors", "14", "price", "discount"},
			expectedRouteKey: fmt.Sprintf("cars:%v:colors:14:price:discount", id),
			expectedIndexKey: fmt.Sprintf("cars:%v", id),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routeKey, indexKey := GetRouteKey(tt.route, tt.id, tt.Args...)
			assert.Equal(t, tt.expectedRouteKey, routeKey)
			assert.Equal(t, tt.expectedIndexKey, indexKey)
		})

	}

}
