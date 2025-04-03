// Package dynamicquerykit provides helper functions to manage:
// Caching with memcached, dynamic filter for queries, standard responses for
// pagination, deleted/updated/created assets and deleted cache
// provides distinct fields struct and an id that can be used for faster lookups
// using the index of each field. Helper methods that generate cache keys based on applied filters
// helper functions that execute database queries for update/delete/update and helper methods for pagination
// both for manual query definition and automatically by wrapping your main query in a CTE.
// database validation methods for user friendly messaged that can be shown as a response
package dynamicquerykit

import "time"

// ErrorResponse standard for errors
type ErrorResponse struct {
	Status  int    `json:"status" xml:"status" yaml:"status" csv:"status"`
	Message string `json:"message" xml:"message" yaml:"message" csv:"message"`
}

// Pagination standard for pagination
type Pagination struct {
	TotalAssets int  `json:"total_assets" xml:"total_assets" yaml:"total_assets" csv:"total_assets"`
	CurrentPage int  `json:"current_page" xml:"current_page" yaml:"current_page" csv:"current_page"`
	TotalPages  int  `json:"total_pages" xml:"total_pages" yaml:"total_pages" csv:"total_pages"`
	NextPage    bool `json:"next_page" xml:"next_page" yaml:"next_page" csv:"next_page"`
	Limit       int  `json:"limit" xml:"limit" yaml:"limit" csv:"limit"`
}

// Filters standard for allowed filters
type Filters struct {
	Name     string `json:"name" xml:"name" yaml:"name" csv:"name"`
	Operator string `json:"operator" xml:"operator" yaml:"operator" csv:"operator"`
	DbField  string `json:"db_field" xml:"db_field" yaml:"db_field" csv:"db_field"`
	FieldID  string `json:"field_id" xml:"field_id" yaml:"field_id" csv:"field_id"`
}

// DeletedCacheResponse standard response for routes that delete cache
type DeletedCacheResponse struct {
	Status      int `json:"status" xml:"status" yaml:"status" csv:"status"`
	KeysFlushed int `json:"keys_flushed" xml:"keys_flushed" yaml:"keys_flushed" csv:"keys_flushed"`
}

// CreatedAssetResponse generic response when creating a new asset
type CreatedAssetResponse struct {
	Status    int       `json:"status" xml:"status" yaml:"status" csv:"status"`
	AssetID   *int64    `json:"asset_id" xml:"asset_id" yaml:"asset_id" csv:"asset_id"`
	CreatedAt time.Time `json:"date" xml:"date" yaml:"date" csv:"date"`
}

// DistinctFieldNamesResponse response for when selecting a field from a route
type DistinctFieldNamesResponse struct {
	Status int                  `json:"status" xml:"status" yaml:"status" csv:"status"`
	Total  int                  `json:"total" xml:"total" yaml:"total" csv:"total"`
	Data   []DistinctFieldNames `json:"data" xml:"data" yaml:"data" csv:"data"`
}

// DistinctFieldNames is the basic component of DistinctFieldNamesResponse
type DistinctFieldNames struct {
	ID    string `json:"id" xml:"id" yaml:"id" csv:"id"`
	Name  string `json:"name" xml:"name" yaml:"name" csv:"name"`
	Count int    `json:"count" xml:"count" yaml:"count" csv:"count"`
}

// IDs generic struct for getting the IDs of assets
type IDs struct {
	ID int `json:"id" xml:"id" yaml:"id" csv:"id"`
}
