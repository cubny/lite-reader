package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// errorType is used for error codes
type errorType int

var (
	// List of API error types.
	errInternalError errorType = 500
	errBadRequest    errorType = 400
	errInvalidParams errorType = 422
	errNotFound      errorType = 404
)

// JSONError is used to return http errors encoded in json
type JSONError struct {
	// Code of the error
	Code int `json:"code"`
	// Details of the error
	Details string `json:"details"`
}

func newJSONError(errorType errorType, details string) JSONError {
	e := JSONError{Code: int(errorType)}

	switch errorType {
	case errInternalError:
		e.Details = "Internal error"
	case errBadRequest:
		e.Details = "Bad Request"
	case errInvalidParams:
		e.Details = "Invalid params"
	case errNotFound:
		e.Details = "Not found"
	default:
		e.Code = 100999
		e.Details = "Unknown error"
	}

	if details != "" {
		e.Details = fmt.Sprintf("%s - %s", e.Details, details)
	}

	return e
}

func (e JSONError) write(w http.ResponseWriter, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := struct {
		Err JSONError `json:"error"`
	}{
		Err: e,
	}

	return json.NewEncoder(w).Encode(resp)
}

// InternalError writes the error details in json with the provided details
func InternalError(w http.ResponseWriter, details string) error {
	return newJSONError(errInternalError, details).write(w, http.StatusInternalServerError)
}

// BadRequest writes the BadRequest error details in json with the provided details
func BadRequest(w http.ResponseWriter, details string) error {
	return newJSONError(errBadRequest, details).write(w, http.StatusBadRequest)
}

// InvalidParams writes the UnprocessableEntity error details in json with the provided details
func InvalidParams(w http.ResponseWriter, details string) error {
	return newJSONError(errInvalidParams, details).write(w, http.StatusUnprocessableEntity)
}

// NotFound writes the NotFound error details in json with the provided details
func NotFound(w http.ResponseWriter, details string) error {
	return newJSONError(errNotFound, details).write(w, http.StatusNotFound)
}
