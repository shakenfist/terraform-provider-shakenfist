package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/shakenfist/client-go"
)

func resourceKey() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace of the key",
				ForceNew:    true,
			},
			"keyname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the key",
				ForceNew:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The access key",
				Sensitive:   true,
			},
		},
		Create: resourceCreateKey,
		Read:   resourceReadKey,
		Delete: resourceDeleteKey,
		Exists: resourceExistsKey,
		Update: resourceUpdateKey,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateKey(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	err := apiClient.CreateNamespaceKey(
		d.Get("namespace").(string),
		d.Get("keyname").(string),
		d.Get("key").(string),
	)
	if err != nil {
		return fmt.Errorf("Unable to create key: %v", err)
	}

	d.SetId(d.Get("keyname").(string))

	return nil
}

func resourceReadKey(d *schema.ResourceData, m interface{}) error {
	// No data can be read from a namespace access key.
	return nil
}

func resourceDeleteKey(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	err := apiClient.DeleteNamespaceKey(d.Get("namespace").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Unable to delete namespace key: %v", err)
	}
	d.SetId("")
	return nil
}

func resourceExistsKey(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	keynames, err := apiClient.GetNamespaceKeys(d.Get("namespace").(string))
	if err != nil {
		return false, fmt.Errorf("Unable to retrieve namespace keys: %v", err)
	}

	for _, n := range keynames {
		if n == d.Id() {
			return true, nil
		}
	}

	return false, nil
}

func resourceUpdateKey(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	if d.HasChange("key") {
		err := apiClient.UpdateNamespaceKey(
			d.Get("namespace").(string),
			d.Get("keyname").(string),
			d.Get("key").(string),
		)
		if err != nil {
			return fmt.Errorf(
				"UpdateNamespace: cannot update namespace key: %v", err)
		}
	}

	return nil
}
