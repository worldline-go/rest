package rest

import (
	"encoding/json"
	"fmt"
	"io"
)

func BindJSON[T any](data io.Reader, result T) error {
	return json.NewDecoder(data).Decode(result)
}

func BindJSONList[T any](data io.Reader, result *[]T) error {
	// Read all data into a buffer to allow multiple parsing attempts
	bodyBytes, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("read request body: %w", err)
	}

	// First try to decode as an array
	if err := json.Unmarshal(bodyBytes, result); err == nil {
		return nil
	}

	// If that fails, try to decode as a single object
	var single T
	if err := json.Unmarshal(bodyBytes, &single); err != nil {
		return err
	}

	*result = []T{single}

	return nil
}
