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
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString},
			},
		},
		Create: resourceCreateNamespace,
		Read:   resourceReadNamespace,
		Delete: resourceDeleteNamespace,
		Exists: resourceExistsNamespace,
		Update: resourceUpdateNamespace,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateNamespace(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	namespace := d.Get("name").(string)

	if err := apiClient.CreateNamespace(namespace); err != nil {
		return fmt.Errorf("Unable to create namespace: %v", err)
	}

	// Set metadata on namespace
	for k, v := range d.Get("metadata").(map[string]interface{}) {
		val, ok := v.(string)
		if !ok {
			return fmt.Errorf("Tag value is not a string")
		}

		err := apiClient.SetMetadata(client.TypeNamespace, namespace, k, val)
		if err != nil {
			return fmt.Errorf("CreateNamespace cannot store metadata: %v", err)
		}
	}

	d.SetId(namespace)

	return nil
}

func resourceReadNamespace(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	// Retrieve metadata
	metadata, err := apiClient.GetMetadata(client.TypeNamespace, d.Id())
	if err != nil {
		return fmt.Errorf("ReadNamespace unable to retrieve metadata: %v", err)
	}
	if err := d.Set("metadata", metadata); err != nil {
		return fmt.Errorf("Namespace Metadata cannot be set: %v", err)
	}

	return nil
}

func resourceDeleteNamespace(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	err := apiClient.DeleteNamespace(d.Id())
	if err != nil {
		return fmt.Errorf("Unable to delete namespace: %v", err)
	}
	d.SetId("")
	return nil
}

func resourceExistsNamespace(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	namespaces, err := apiClient.GetNamespaces()
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

func resourceUpdateNamespace(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	if d.HasChange("metadata") {
		if err := updateMetadata(client.TypeNamespace, d, apiClient); err != nil {
			return fmt.Errorf("UpdateNamespace error: %v", err)
		}
	}

	return nil
}
