package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
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

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := errors.Cause(err).(*Error); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err := Respond(w, er, webErr.Status); err != nil {
			return err
		}

		return nil
	}

	// If not, the handler sent any arbitrary error value so use 500
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(w, er, http.StatusInternalServerError); err != nil {
		return err
	}

	return nil
}
