// With many thanks to the example code from
// https://github.com/spaceapegames/terraform-provider-example
package provider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/shakenfist/client-go"
)

func validateNetblock(v interface{}, k string) ([]string, []error) {
	var errs []error
	var warns []string

	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected netblock to be a string"))
		return warns, errs
	}

	netblock := regexp.MustCompile(
		`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}/[0-9]{1,3}$`)
	if !netblock.Match([]byte(value)) {
		errs = append(errs,
			fmt.Errorf("Netblock must be IPv4 CIDR. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceNetwork() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the network",
				ForceNew:    true,
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the network",
			},
			"netblock": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The CIDR IP range of the network",
				ForceNew:     true,
				ValidateFunc: validateNetblock,
			},
			"provide_dhcp": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Should DHCP services exist on the network?",
				ForceNew:    true,
			},
			"provide_nat": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Should NAT services exist on the network?",
				ForceNew:    true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString},
			},
		},
		Create: resourceCreateNetwork,
		Read:   resourceReadNetwork,
		Delete: resourceDeleteNetwork,
		Exists: resourceExistsNetwork,
		Update: resourceUpdateNetwork,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateNetwork(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	network, err := apiClient.CreateNetwork(
		d.Get("netblock").(string), d.Get("provide_dhcp").(bool),
		d.Get("provide_nat").(bool), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Unable to create network: %v", err)
	}

	if err := d.Set("uuid", network.UUID); err != nil {
		return fmt.Errorf("UUID cannot be set: %v", err)
	}

	d.SetId(network.UUID)

	// Set metadata on the network
	for k, v := range d.Get("metadata").(map[string]interface{}) {
		val, ok := v.(string)
		if !ok {
			return fmt.Errorf("Tag value is not a string")
		}

		err := apiClient.SetMetadata(client.TypeNetwork, network.UUID, k, val)
		if err != nil {
			return fmt.Errorf("CreateNetwork cannot store metadata: %v", err)
		}
	}

	return nil
}

func resourceReadNetwork(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	network, err := apiClient.GetNetwork(d.Id())
	if err != nil {
		return fmt.Errorf("Unable to retrieve network: %v", err)
	}

	if err := d.Set("uuid", network.UUID); err != nil {
		return fmt.Errorf("Network UUID cannot be set: %v", err)
	}
	if err := d.Set("name", network.Name); err != nil {
		return fmt.Errorf("Network Name cannot be set: %v", err)
	}
	if err := d.Set("netblock", network.NetBlock); err != nil {
		return fmt.Errorf("Network NetBlock cannot be set: %v", err)
	}
	if err := d.Set("provide_dhcp", network.ProvideDHCP); err != nil {
		return fmt.Errorf("Network ProvideDHCP flag cannot be set: %v", err)
	}
	if err := d.Set("provide_nat", network.ProvideNAT); err != nil {
		return fmt.Errorf("Network ProvideNAT flag cannot be set: %v", err)
	}

	// Retrieve metadata
	metadata, err := apiClient.GetNetworkMetadata(d.Id())
	if err != nil {
		return fmt.Errorf("ReadInstance unable to retrieve metadata: %v", err)
	}
	if err := d.Set("metadata", metadata); err != nil {
		return fmt.Errorf("Network Metadata cannot be set: %v", err)
	}

	return nil
}

func resourceDeleteNetwork(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	err := apiClient.DeleteNetwork(d.Id())
	if err != nil {
		return fmt.Errorf("Unable to delete network: %v", err)
	}
	d.SetId("")
	return nil
}

func resourceExistsNetwork(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	n, err := apiClient.GetNetwork(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, fmt.Errorf("Unable to check network existence: %v", err)
		}
	}

	if n.State == "deleted" {
		return false, nil
	}

	return true, nil
}

func resourceUpdateNetwork(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	if d.HasChange("metadata") {
		if err := updateMetadata(client.TypeNetwork, d, apiClient); err != nil {
			return fmt.Errorf("UpdateNetwork error: %v", err)
		}
	}

	return nil
}
