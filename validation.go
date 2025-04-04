package dynamicquerykit

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
)

// IsFieldJSONTag checks if the string provided matches any json tag of a struct
// it requires the struct provided has json tags specified
func IsFieldJSONTag(dataStruct any, strField string) bool {
	method := "IsFieldJSONTag"

	res := false
	reflectionStruct := reflect.TypeOf(dataStruct)

	for i := range reflectionStruct.NumField() {
		slog.LogAttrs(context.Background(), slog.LevelDebug, "struct reflection", slog.String("method", method), slog.Any("field", reflectionStruct.Field(i)), slog.String("string", strField))
		jsonTag := strings.Split(reflectionStruct.Field(i).Tag.Get("json"), ",")[0]
		if strField == jsonTag {
			slog.LogAttrs(context.Background(), slog.LevelDebug, "string match", slog.String("method", method), slog.String(strField, jsonTag))
			res = true
		}
	}
	return res
}

// IsFieldFilter returns true if a string field matches the name of any filter
// in the provided filter slice
func IsFieldFilter(filters []Filters, field string) (bool, Filters) {
	for _, filter := range filters {
		if filter.Name == field || filter.DbField == field {
			return true, filter
		}
	}
	return false, Filters{}
}

// OrderValidation provide a struct, the order by string and d.irection and default order by string
func OrderValidation(orderByStr string, orderDirectionStr string, filters []Filters) string {
	method := "OrderValidation"
	orderBy := strings.ToLower(orderByStr)
	orderDirection := strings.ToUpper(orderDirectionStr)

	slog.LogAttrs(context.Background(), slog.LevelDebug, "validating order field", slog.String("method", method), slog.String("order by", orderBy), slog.String("order direction", orderDirection))
	if orderDirection != "DESC" && orderDirection != "ASC" {
		orderDirection = "ASC"
	}

	result, filter := IsFieldFilter(filters, orderBy)
	if !result || orderBy == "" {
		slog.LogAttrs(context.Background(), slog.LevelWarn, "order by that does not exist as a filter was provided, using first field instead", slog.String("order by", orderBy), slog.String("order direction", orderDirection), slog.String("filter", filters[0].DbField))
		orderBy = filters[0].DbField
	} else {
		orderBy = filter.DbField
	}
	return fmt.Sprintf("%s %s", orderBy, orderDirection)
}

// DatabaseValidation checks a database error and returns an appropriate
// error message and status code that can be directly used in the response.
// The underline error messaages is always logged.
func DatabaseValidation(err error) (int, error) {
	slog.LogAttrs(context.Background(), slog.LevelError, "database error",
		slog.String("error", err.Error()),
	)
	if strings.Contains(err.Error(), "1062") {
		return http.StatusConflict, fmt.Errorf("Asset already exists")
	}
	if strings.Contains(err.Error(), "1452") {
		return http.StatusBadRequest, fmt.Errorf("Invalid related resource")
	}
	if strings.Contains(err.Error(), "1406") {
		return http.StatusBadRequest, fmt.Errorf("Data too long for column")
	}
	if strings.Contains(err.Error(), "1048") {
		return http.StatusBadRequest, fmt.Errorf("Column cannot be null")
	}
	if strings.Contains(err.Error(), "no rows in result set") {
		return http.StatusNotFound, fmt.Errorf("No data for specified request")
	}

	return http.StatusInternalServerError, fmt.Errorf("An unexpected error occurred")
}
