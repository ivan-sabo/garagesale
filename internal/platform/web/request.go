package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Decode looks for a JSON document in the request body and unmarshals it into val
func Decode(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return fmt.Errorf("decoding request body: %w", err)
	}

	return nil
}
