package multiverse

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Provider ...
func Provider() *schema.Provider {
	binaryName := filepath.Base(os.Args[0])

	resourceName := strings.TrimPrefix(binaryName, "terraform-provider-")
	log.Printf("multiverse Found resource type: %s\n", resourceName)
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			resourceName: resourceCustom(),
		},
	}
	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {

		return nil, nil
	}
}
