package provider

import (
	"fmt"
)

func compareMetadata(correctMeta, serverMeta map[string]string) error {
	// Test correct metadata against the actual instance
	for k, v := range correctMeta {
		actualV, ok := serverMeta[k]

		if !ok {
			return fmt.Errorf("Created server object missing key: %s", k)
		}

		if actualV != v {
			return fmt.Errorf(
				"Metadata key %s has wrong value: %s (should be %s)",
				k, actualV, v)
		}
	}

	return nil
}
