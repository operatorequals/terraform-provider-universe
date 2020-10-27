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

const (
	DefaultProviderName = "multiverse"
	EnvProviderNameVar  = "TERRAFORM_MULTIVERSE_PROVIDERNAME"
)

// getProviderNameFromBinaryOrEnvironment
// if the user has explicitly specified provider name in env var use it,
// otherwise look for it in th binary name after "terraform-provider-"
// but if debugging or testing the binary name is junk e.g. 'debug.test'
// so provide a default.
func getProviderNameFromBinaryOrEnvironment() (name string) {
	const TFP string = "terraform-provider-"

	name, ok := os.LookupEnv(EnvProviderNameVar)
	if ok {
		return // env var overrides binary name
	}
	binaryName := filepath.Base(os.Args[0])
	if strings.HasPrefix(binaryName, TFP) {
		name = strings.TrimPrefix(binaryName, TFP)
		return
	}
	name = DefaultProviderName
	return
}

// getResourceTypeNamesFromEnvironment
// Assuming the environment has a variable TERRAFORM_MULTIVERSE_RESOURCETYPES containing a
// whitespace-separated list of resource names.
// Return a []string of the names plus "multiverse"
func getResourceTypeNamesFromEnvironment(providerName string) (result map[string]bool) {
	result = map[string]bool{providerName: true}
	prefix := providerName + "_"

	resourceTypesVarName := "TERRAFORM_" + strings.ToUpper(providerName) + "_RESOURCETYPES"
	resourceTypeNames, ok := os.LookupEnv(resourceTypesVarName)
	if !ok {
		return
	}
	f := strings.Fields(resourceTypeNames)

	for _, x := range f {
		rtn := x
		if !strings.HasPrefix(x, prefix) { // Enforce rule that resource type names must be providername '_' resoyrcetypename
			rtn = prefix + x
		}
		result[rtn] = true
	}
	return
}

func getResourceMap(providerName string) (result map[string]*schema.Resource) {
	result = make(map[string]*schema.Resource)
	for resourceName := range getResourceTypeNamesFromEnvironment(providerName) {
		result[resourceName] = resourceCustom()
	}
	log.Printf("resourceMap is: %#v\n", result)
	return
}

// Provider ...
func Provider() *schema.Provider {
	// Get the provider name to use
	providerName := getProviderNameFromBinaryOrEnvironment()
	log.Printf("multiverse provider name is: %s\n", providerName)

	// Get the resource names
	resourceMap := getResourceMap(providerName)
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
		configurationData := map[string]interface{}{}
		for _, key := range []string{"id_key", "executor", "script", "environment", "computed", "javascript"} {
			val, ok := d.GetOk(key)
			if !ok {
				continue
			}
			configurationData[key] = val
		}
		// Just check the environment is a map
		e, ok := d.GetOk("environment")
		if ok {
			if _, ok = e.(map[string]interface{}); !ok {
				return nil, fmt.Errorf("as expecting map[string]interface{} in 'environment', got %#v", e)
			}
		}
		return configurationData, nil
	}

	return p
}
