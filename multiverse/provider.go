package multiverse

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// getResourceNamesFromEnvironment
// assuming the environment has a variable MULTIVERSE_RESOURCENAMES containign a
// comma-seperated list of resource names.
// Return a []string of the names plus "multiverse"
func getResourceNamesFromEnvironment(providerName string) (result []string) {
	result = make([]string, 10)
	resourceNamesName := "TERRAFORM_" + strings.ToUpper(providerName) + "_RESOURCENAMES"
	resourceNames, ok := os.LookupEnv(resourceNamesName)
	if ok {
		result = strings.Split(resourceNames, ",")
	}
	result = append(result, providerName)
	return
}

func getResourceMap(providerName string) (result map[string]*schema.Resource) {
	result = make(map[string]*schema.Resource)
	for _, resourceName := range getResourceNamesFromEnvironment(providerName) {
		result[resourceName] = resourceCustom()
	}
	return
}

// Provider ...
func Provider() *schema.Provider {
	// Get the provider name to use
	binaryName := filepath.Base(os.Args[0])
	providerName := strings.TrimPrefix(binaryName, "terraform-provider-")

	// Get the resource names
	resourceMap := getResourceMap(providerName)
	log.Printf("multiverse provider name : %s\n", providerName)
	for n := range resourceMap {
		log.Printf("provider %s has resource %s\n", providerName, n)
	}

	p := &schema.Provider{
		ResourcesMap: resourceMap,
		Schema: map[string]*schema.Schema{
			"id_key": {
				Description: "The name of the key which holds the unique identifier of the resource. e.g. 'id'",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"executor": {
				Description: "The name of the program to run. e.g. python",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"script": {
				Description: "The path to the script passed as the first argument to 'executor'.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"environment": {
				Description: "The configuration passed as environment variables to the provider script.",
				Optional:    true,
				Type:        schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"computed": {
				Description:  "A list of fields (in JSON format) returned from the executor script which are computed dynamically.",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsJSON,
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
