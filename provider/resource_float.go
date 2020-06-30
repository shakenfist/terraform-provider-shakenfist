// With many thanks to the example code from
// https://github.com/spaceapegames/terraform-provider-example
package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/shakenfist/client-go"
)

func resourceFloat() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"interface": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of Interface",
				ForceNew:    true,
			},
			"ipv4": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IPv4 Address",
			},
		},
		Create: resourceCreateFloat,
		Read:   resourceReadFloat,
		Delete: resourceDeleteFloat,
		Exists: resourceExistsFloat,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateFloat(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	uuid := d.Get("interface").(string)

	err := apiClient.FloatInterface(uuid)
	if err != nil {
		return fmt.Errorf("Unable to float interface: %v", err)
	}

	d.Set("uuid", uuid)
	d.SetId(uuid)

	if err := resourceReadFloat(d, m); err != nil {
		return fmt.Errorf("CreateFloat: %v", err)
	}

	return nil
}

func resourceReadFloat(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	iface, err := apiClient.GetInterface(d.Id())
	if err != nil {
		return fmt.Errorf("Unable to retrieve network: %v", err)
	}

	if iface.Floating == "" {
		return fmt.Errorf("Interface does not have a floating IP")
	}

	d.Set("ipv4", iface.Floating)
	return nil
}

func resourceDeleteFloat(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	err := apiClient.DefloatInterface(d.Id())
	if err != nil {
		return fmt.Errorf("Unable to delete network interface: %v", err)
	}
	d.SetId("")
	return nil
}

func resourceExistsFloat(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	iface, err := apiClient.GetInterface(d.Id())
	if err != nil {
		return false, fmt.Errorf("Unable to retrieve network interface: %v", err)
	}

	if iface.Floating == "" {
		return false, nil
	}

	return true, nil
}
