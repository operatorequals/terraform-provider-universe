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
		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The configuration passed as the last argument to the provider script.",
			},
		},
	}
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		conf, ok := d.GetOk("configuration")
		if !ok {
			return "", nil
		}
		config := conf.(string)
		return config, nil
	}

	return p
}
