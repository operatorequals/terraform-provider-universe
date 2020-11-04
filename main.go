package main

import (
	"github.com/birchb1024/terraform-provider-universe/universe"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: providerProvider,
	})
}

func providerProvider() *schema.Provider {
	return universe.Provider()
}
