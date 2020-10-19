package multiverse

import (
	"fmt"
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
			"id_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the key which holds the unique identifier of the resource. e.g. 'id'",
			},
			"executor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the program to run. e.g. python",
			},
			"script": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to the script passed as the first argument to 'executor'.",
			},
			"environment": {
				Optional:    true,
				Description: "The configuration passed as environment variables to the provider script.",
				Type:        schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		result := map[string]interface{}{}
		e, ok := d.GetOk("environment")
		if !ok {
			return "", nil
		}
		_, ok = e.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("as expecting map[string]interface{} in 'environment', got %v", e)
		}
		for _, key := range []string{"id_key", "executor", "script", "environment"} {
			val, ok := d.GetOk(key)
			if !ok {
				continue
			}
			result[key] = val
		}
		return result, nil
	}

	return p
}
