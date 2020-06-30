package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/shakenfist/client-go"
)

func resourceNamespace() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the namespace",
				ForceNew:    true,
			},
			"keyname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the key",
				ForceNew:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The key used for authentication",
				ForceNew:    true,
			},
		},
		Create: resourceCreateNamespace,
		Read:   resourceReadNamespace,
		Delete: resourceDeleteNamespace,
		Exists: resourceExistsNamespace,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateNamespace(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	err := apiClient.CreateNameSpace(
		d.Get("name").(string),
		d.Get("keyname").(string),
		d.Get("key").(string),
	)
	if err != nil {
		return fmt.Errorf("Unable to create namespace: %v", err)
	}

	d.SetId(d.Get("name").(string))

	return nil
}

func resourceReadNamespace(d *schema.ResourceData, m interface{}) error {

	// A Namespace only has keynames and keys.
	// Neither are relevant at the moment to Terraform

	return nil
}

func resourceDeleteNamespace(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	err := apiClient.DeleteNameSpace(d.Id())
	if err != nil {
		return fmt.Errorf("Unable to delete namespace: %v", err)
	}
	d.SetId("")
	return nil
}

func resourceExistsNamespace(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	namespaces, err := apiClient.GetNameSpaces()
	if err != nil {
		return false, fmt.Errorf("Unable to retrieve namespaces: %v", err)
	}

	for _, n := range namespaces {
		if n == d.Id() {
			return true, nil
		}
	}

	return false, nil
}
