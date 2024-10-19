package nubio

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadJSONFile(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read whole file: %w", err)
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return fmt.Errorf("unmarshal JSON: %w", err)
	}
	return nil
}
