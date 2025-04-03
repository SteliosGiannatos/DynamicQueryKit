package dynamicquerykit

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

// GetString gets env variable of a string and if missing sets a default value
func GetString(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

// GetTime Gets env variable of a key, makes it into time.duration else returns the fallback
func GetTime(key, fallback string) time.Duration {
	val, ok := os.LookupEnv(key)
	f, err := time.ParseDuration(fallback)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "you provided an incorrect fallback duration\n", slog.String("key", fallback), slog.Any("error", err))
		panic(fmt.Sprintf("you provided an incorrect fallback duration for env key: (%v), error: (%v)\n", fallback, err))

	}
	if !ok {
		return f
	}
	dur, err := time.ParseDuration(val)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelWarn, "you provided an incorrect duration key, using fallback\n", slog.String("key", key), slog.String("value", val))
		return f

	}
	return dur
}

// GetInt gets env variable of a int and if missing sets a default value
func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}

// GetBool gets env variable of a bool and if missing sets a default value
func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}

// GetDates gets a map of named string dates and returns them as time.Time
// the dates must be using the time.DateOnly format
// const time.DateOnly untyped string = "2006-01-02"
func GetDates(dateMap map[string]string) map[string]time.Time {
	dates := make(map[string]time.Time)
	parseDate := func(date_str string) time.Time {
		parsed, err := time.Parse(time.DateOnly, date_str)
		if err != nil {
			return time.Time{}
		}
		return parsed
	}

	for k, v := range dateMap {
		date := parseDate(v)
		if !date.IsZero() {
			dates[k] = date
		}
	}

	return dates
}
