package main

import (
	"github.com/birchb1024/terraform-provider-multiverse/multiverse"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return multiverse.Provider()
		},
	})
}
