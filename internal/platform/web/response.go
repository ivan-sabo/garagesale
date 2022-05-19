package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Respond marshals a value to JSON and sends it to the client
func Respond(w http.ResponseWriter, value interface{}, statusCode int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshaling value to json: %w", err)
	}

	w.Header().Set("content-type", "application/json; charset= utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("writing to client: %w", err)
	}

	return nil
}

// ResponError knows how to handle errors going out to the client
func ResponError(w http.ResponseWriter, err error) error {
	if webErr, ok := err.(*Error); ok {
		resp := ErrorResponse{
			Error: webErr.Err.Error(),
		}

		return Respond(w, resp, webErr.Status)
	}

	resp := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	return Respond(w, resp, http.StatusInternalServerError)
}
