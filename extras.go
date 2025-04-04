package dynamicquerykit

import (
	sq "github.com/Masterminds/squirrel"
)

// GetPaginationQuery provides a query that can be used for pagination.
// it wraps the provided query as a subquery.
func GetPaginationQuery(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Prefix(`SELECT COUNT(*) AS total_rows FROM (`).Suffix(`) AS grouped_results;`)
}
