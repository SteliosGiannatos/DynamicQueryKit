package dynamicquerykit

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AcceptedEncoding returns a default encoding if the one specified does not match the specification
func AcceptedEncoding(accept string) string {
	accept = strings.ToLower(strings.TrimSpace(accept))
	supportedEncodings := map[string]bool{
		"application/json": true,
		"application/xml":  true,
	}

	if supportedEncodings[accept] {
		return accept
	}

	return "application/json"
}

// DataEncode encodes structs to json and xml
// provides an easily scalable way to support extra encodings accross the api
// Returns the Data encoded, returns the Content-Type of the encoded data and error
func DataEncode(Accept string, Data any) ([]byte, string, error) {
	var (
		err         error
		encodedData []byte
		contentType string
	)

	switch strings.ToLower(strings.TrimSpace(Accept)) {
	case "application/json":
		contentType = "application/json"
		encodedData, err = json.Marshal(Data)
	case "application/xml":
		contentType = "application/xml"
		encodedData, err = xml.Marshal(Data)
	default:
		contentType = "application/json"
		encodedData, err = json.Marshal(Data)
	}

	return encodedData, contentType, err
}

// DecodeBody gets the body of a request and parses it
func DecodeBody(contentType string, body io.Reader, data any) (int, error) {
	var (
		err                error
		format             string
		msg                string
		statusCode         int
		jsonDecoder        *json.Decoder
		xmlDecoder         *xml.Decoder
		syntaxErrorJSON    *json.SyntaxError
		syntaxErrorXML     *xml.SyntaxError
		unmarshalTypeError *json.UnmarshalTypeError
	)
	switch contentType {
	case "application/json":
		format = "JSON"
		jsonDecoder = json.NewDecoder(body)
		jsonDecoder.DisallowUnknownFields()
		err = jsonDecoder.Decode(data)

	case "application/xml":
		format = "XML"
		xmlDecoder = xml.NewDecoder(body)
		err = xmlDecoder.Decode(data)

	default:
		return http.StatusBadRequest, fmt.Errorf("failed to decode body for specified Content-Type")
	}

	if err != nil {
		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case format == "JSON" && errors.As(err, &syntaxErrorJSON):
			msg = fmt.Sprintf("Request body contains badly-formed %s (at position %d)", format, syntaxErrorJSON.Offset)
			statusCode = http.StatusBadRequest

		case format == "XML" && errors.As(err, &syntaxErrorXML):
			msg = fmt.Sprintf("Request body contains badly-formed %s (at line %d)", format, syntaxErrorXML.Line)
			statusCode = http.StatusBadRequest

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg = fmt.Sprintf("Request body contains badly-formed %s", format)
			statusCode = http.StatusBadRequest

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case format == "JSON" && errors.As(err, &unmarshalTypeError):
			msg = fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			statusCode = http.StatusBadRequest

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.

		//This does not apply to xml, since there is no xmlDecoder.DisallowUnknownFields..
		//but it makes it easier to extend in the future
		case strings.HasPrefix(err.Error(), fmt.Sprintf("%s: unknown field ", strings.ToLower(format))):
			fieldName := strings.TrimPrefix(err.Error(), fmt.Sprintf("%s: unknown field ", strings.ToLower(format)))
			msg = fmt.Sprintf("Request body contains unknown field %s", fieldName)
			statusCode = http.StatusBadRequest

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg = "Request body must not be empty"
			statusCode = http.StatusBadRequest

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg = "Request body must not be larger than 1MB"
			statusCode = http.StatusRequestEntityTooLarge

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			msg = err.Error()
			statusCode = http.StatusInternalServerError
		}
		return statusCode, fmt.Errorf("%s", msg)
	}
	switch format {
	case "JSON":
		err = jsonDecoder.Decode(&struct{}{})
		if !errors.Is(err, io.EOF) {
			return http.StatusBadRequest, fmt.Errorf("Request body must only contain a single JSON object")
		}
	case "XML":
		err = xmlDecoder.Decode(&struct{}{})
		if !errors.Is(err, io.EOF) {
			return http.StatusBadRequest, fmt.Errorf("Request body must only contain a single JSON object")
		}
	}

	return http.StatusOK, nil
}
