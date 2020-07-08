package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/shakenfist/client-go"
)

// UpdateMetadata compares the old and new metadata, then updates the
// resource metadata on Shaken Fist.
func updateMetadata(
	resType client.ResourceType,
	d *schema.ResourceData,
	apiClient *client.Client) error {

	var err error

	// Retrieve metadata changes
	o, n := d.GetChange("metadata")
	oldMeta := o.(map[string]interface{})
	newMeta := n.(map[string]interface{})

	// Update new and changed metadata
	for key, newVal := range newMeta {
		if oldVal, exists := oldMeta[key]; exists {
			if oldVal != newVal {
				// Old key, value changing
				err = apiClient.SetMetadata(resType, d.Id(), key, newVal.(string))
			}
		} else {
			// New key
			err = apiClient.SetMetadata(resType, d.Id(), key, newVal.(string))
		}

		if err != nil {
			return fmt.Errorf("Unable to change metadata: %v", err)
		}
	}

	// Find deleted metadata keys
	for key := range oldMeta {
		if _, exists := newMeta[key]; !exists {
			err = apiClient.DeleteMetadata(resType, d.Id(), key)
			if err != nil {
				return fmt.Errorf("Unable to delete metadata key: %v", err)
			}
		}
	}

	return nil
}
