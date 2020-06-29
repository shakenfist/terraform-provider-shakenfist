// With many thanks to the example code from
// https://github.com/spaceapegames/terraform-provider-example
package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	client "github.com/shakenfist/client-go"
)

// Provider is the terraform provider interface
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SHAKENFIST_SERVER", ""),
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SHAKENFIST_PORT", ""),
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SHAKENFIST_NAMESPACE", ""),
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SHAKENFIST_KEY", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"shakenfist_network":  resourceNetwork(),
			"shakenfist_instance": resourceInstance(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	address := d.Get("address").(string)
	port := d.Get("port").(int)
	namespace := d.Get("namespace").(string)
	key := d.Get("key").(string)

	return client.NewClient(address, port, namespace, key), nil
}
