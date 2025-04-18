package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var ErrDecode = errors.New("failed to decode data")

func BindJSONList[T any](data io.Reader, result *[]T) error {
	// Read all data into a buffer to allow multiple parsing attempts
	bodyBytes, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	// First try to decode as an array
	if err := json.Unmarshal(bodyBytes, result); err == nil {
		return nil
	}

	// If that fails, try to decode as a single object
	var single T
	if err := json.Unmarshal(bodyBytes, &single); err == nil {
		*result = []T{single}

		return nil
	}

	return ErrDecode
}
