package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/operatorequals/terraform-provider-universe/universe"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: providerProvider,
	})
}

func providerProvider() *schema.Provider {
	return universe.Provider()
}
