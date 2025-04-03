package dynamicquerykit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	v := "SHELL"
	t.Setenv(v, "/bin/bash")

	tests := []struct {
		name     string
		key      string
		fallback string
		expected string
	}{
		{
			name:     "returns existing env variable",
			fallback: "/bin/zsh",
			key:      v,
			expected: "/bin/bash",
		},
		{
			name:     "returns fallback if env variable not defined",
			fallback: "/bin/zsh",
			key:      "UndefinedEnvVariable",
			expected: "/bin/zsh",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := GetString(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, val)
		})

	}

}

func TestGetTime(t *testing.T) {
	v := "FIVE"
	t.Setenv(v, "5m")

	tests := []struct {
		name     string
		key      string
		fallback time.Duration
		expected time.Duration
	}{
		{
			name:     "returns existing env variable",
			fallback: time.Duration(5 * time.Hour),
			key:      v,
			expected: time.Duration(5 * time.Minute),
		},
		{
			name:     "returns fallback if env variable not defined",
			fallback: time.Duration(8 * time.Hour),
			key:      "UndefinedEnvVariable",
			expected: time.Duration(8 * time.Hour),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := GetTime(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, val)
		})

	}

}

func TestGetInt(t *testing.T) {
	v := "FIVE"
	t.Setenv(v, "5")

	tests := []struct {
		name     string
		key      string
		fallback int
		expected int
	}{
		{
			name:     "returns existing env variable",
			fallback: 8,
			key:      v,
			expected: 5,
		},
		{
			name:     "returns fallback if env variable not defined",
			fallback: 8,
			key:      "UndefinedEnvVariable",
			expected: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := GetInt(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, val)
		})

	}

}

func TestGetBool(t *testing.T) {
	v := "GetTrue"
	t.Setenv(v, "true")

	tests := []struct {
		name     string
		key      string
		fallback bool
		expected bool
	}{
		{
			name:     "returns existing env variable",
			fallback: false,
			key:      v,
			expected: true,
		},
		{
			name:     "returns fallback if env variable not defined",
			fallback: false,
			key:      "UndefinedEnvVariable",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := GetBool(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, val)
		})

	}

}

func TestGetDate(t *testing.T) {
	v := "GetFirstDayOfYearMonth"
	t.Setenv(v, "2025-01-01")

	tests := []struct {
		name     string
		key      string
		fallback time.Time
		expected time.Time
	}{
		{
			name:     "returns existing env variable",
			fallback: time.Date(2025, 03, 30, 0, 0, 0, 0, time.UTC),
			key:      v,
			expected: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "returns fallback if env variable not defined",
			fallback: time.Date(2025, 03, 30, 0, 0, 0, 0, time.UTC),
			key:      "UndefinedEnvVariable",
			expected: time.Date(2025, 03, 30, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := GetDate(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, val)
		})

	}

}
