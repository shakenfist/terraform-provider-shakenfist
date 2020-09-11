// With many thanks to the example code from
// https://github.com/spaceapegames/terraform-provider-example
package main

import (
	"github.com/shakenfist/terraform-provider-shakenfist/provider"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
