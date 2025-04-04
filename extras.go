package dynamicquerykit

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// GetPaginationQuery provides a query that can be used for pagination.
// it wraps the provided query as a subquery.
func GetPaginationQuery(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Prefix(`SELECT COUNT(*) AS total_rows FROM (`).Suffix(`) AS grouped_results;`)
}

// GetRouteKey returns the Route key which is used in the GetCacheKey method
// and the correct Index Key which is going to be used in the SetKeyIndex Method.
// The Index Key will store in a json list the Route Key
// ie. -> properties,1 will return properties:1 for both Route Key & Index Key
// ie. -> properties,1,google_searches will return properties:1:google_searches for Route Key and properties:1 Index Key
// Provides an easier way to handle cache keys for subresources of dynamic routes ie. properties/{some_id}/google_searches
// When accessing stuff in deeper levels than that it is recommended to reuse the same method
// ,keeping the Route Key and using the first level Index Key
// ie. -> properties/{some_id}/google_searches/{some_other_id} should keep the index key of properties:some_id
// but routeKey would be properties:some_id:google_searches:some_other_id returned by GetRouteKey
func GetRouteKey(route string, assetID *int, args ...string) (string, string) {
	routeKey := route

	var id int
	if assetID != nil {
		id = *assetID
		routeKey += fmt.Sprintf(":%d", id)
	}
	indexKey := routeKey

	for _, arg := range args {
		routeKey += fmt.Sprintf(":%s", arg)
	}

	return routeKey, indexKey
}
