// With many thanks to the example code from
// https://github.com/spaceapegames/terraform-provider-example
package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	client "github.com/shakenfist/client-go"
)

// Provider is the terraform provider interface
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SHAKENFIST_URL", ""),
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
			"shakenfist_namespace": resourceNamespace(),
			"shakenfist_key":       resourceKey(),
			"shakenfist_network":   resourceNetwork(),
			"shakenfist_instance":  resourceInstance(),
			"shakenfist_float":     resourceFloat(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	server_url := d.Get("server_url").(string)
	namespace := d.Get("namespace").(string)
	key := d.Get("key").(string)

	if server_url == "" {
		return nil, fmt.Errorf(
			"Server URL not set, expecting \"http://<server>:<port>\"")
	}
	if namespace == "" {
		return nil, fmt.Errorf("Namespace not set")
	}
	if key == "" {
		return nil, fmt.Errorf("Access key not set")
	}

	return client.NewClient(server_url, namespace, key), nil
}
