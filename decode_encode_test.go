package dqk

import (
	"encoding/xml"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AcceptedEncoding returns a default encoding if the one specified does not match the specification
func TestAcceptedEncoding(t *testing.T) {
	tests := []struct {
		name     string
		accept   string
		expected string
	}{
		{
			name:     "acceptable encoding json",
			accept:   "application/json",
			expected: "application/json",
		},
		{
			name:     "acceptable encoding xml",
			accept:   "application/xml",
			expected: "application/xml",
		},
		{
			name:     "not acceptable encoding",
			accept:   "application/x-protobuf",
			expected: "application/json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AcceptedEncoding(tt.accept)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestDataEncode(t *testing.T) {
	type TestStruct struct {
		ID   int    `json:"id" xml:"id"`
		Name string `json:"name" xml:"name"`
	}

	tests := []struct {
		name                string
		accept              string
		data                TestStruct
		expected            string
		expectedContentType string
	}{
		{
			name:                "encode json",
			accept:              "application/json",
			expectedContentType: "application/json",
			data: TestStruct{
				ID: 1, Name: "hello",
			},
			expected: `{"id": 1, "name": "hello"}`,
		},
		{
			name:                "encode xml",
			accept:              "application/xml",
			expectedContentType: "application/xml",
			data: TestStruct{
				ID: 2, Name: "world",
			},
			expected: `<TestStruct><id>2</id><name>world</name></TestStruct>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, contentType, err := DataEncode(tt.accept, tt.data)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedContentType, contentType)

			if contentType == "application/json" {
				assert.JSONEq(t, tt.expected, string(out))
			} else if contentType == "application/xml" {
				// Use xml.Marshal to handle exact formatting
				expected, _ := xml.Marshal(tt.data)
				assert.Equal(t, string(expected), string(out))
			}
		})
	}
}

func TestDecodeBody(t *testing.T) {
	type TestStruct struct {
		ID   int    `json:"id" xml:"id"`
		Name string `json:"name" xml:"name"`
	}

	tests := []struct {
		name           string
		contentType    string
		body           string
		expectStatus   int
		expectError    bool
		expectedParsed *TestStruct
	}{
		{
			name:           "valid JSON",
			contentType:    "application/json",
			body:           `{"id": 1, "name": "Stelios"}`,
			expectStatus:   http.StatusOK,
			expectError:    false,
			expectedParsed: &TestStruct{ID: 1, Name: "Stelios"},
		},
		{
			name:         "invalid JSON syntax",
			contentType:  "application/json",
			body:         `{"id": 1, "name": "Stelios"`, // missing closing brace
			expectStatus: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name:         "invalid JSON type mismatch",
			contentType:  "application/json",
			body:         `{"id": "oops", "name": "Stelios"}`,
			expectStatus: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name:         "unknown JSON field",
			contentType:  "application/json",
			body:         `{"id": 1, "name": "Stelios", "extra": 123}`,
			expectStatus: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name:         "empty JSON body",
			contentType:  "application/json",
			body:         ``,
			expectStatus: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name:           "valid XML",
			contentType:    "application/xml",
			body:           `<TestStruct><id>2</id><name>Anna</name></TestStruct>`,
			expectStatus:   http.StatusOK,
			expectError:    false,
			expectedParsed: &TestStruct{ID: 2, Name: "Anna"},
		},
		{
			name:         "invalid XML syntax",
			contentType:  "application/xml",
			body:         `<TestStruct><id>2<id><name>Anna</name></TestStruct>`,
			expectStatus: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name:         "unsupported content type",
			contentType:  "text/html",
			body:         `<html></html>`,
			expectStatus: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name:         "multiple JSON objects",
			contentType:  "application/json",
			body:         `{"id":1,"name":"one"}{"id":2,"name":"two"}`,
			expectStatus: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name:         "multiple XML objects",
			contentType:  "application/xml",
			body:         `<TestStruct><id>1</id><name>A</name></TestStruct><TestStruct><id>2</id><name>B</name></TestStruct>`,
			expectStatus: http.StatusBadRequest,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var parsed TestStruct
			status, err := DecodeBody(tt.contentType, strings.NewReader(tt.body), &parsed)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedParsed, &parsed)
			}
			assert.Equal(t, tt.expectStatus, status)
		})
	}
}
