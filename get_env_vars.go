package dqk

import (
	"context"
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
func GetTime(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	dur, err := time.ParseDuration(val)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelWarn, "you provided an incorrect duration key, using fallback\n", slog.String("key", key), slog.String("value", val))
		return fallback

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

// GetDate gets a time.DateOnly string and makes it time.Time
func GetDate(key string, fallback time.Time) time.Time {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	parsed, err := time.Parse(time.DateOnly, val)
	if err != nil {
		return fallback
	}

	return parsed
}
