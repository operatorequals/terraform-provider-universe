package multiverse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"os"
	"os/exec"
)

func resourceCustom() *schema.Resource {
	return &schema.Resource{
		Create: onCreate,
		Read:   onRead,
		Update: onUpdate,
		Delete: onDelete,
		Exists: onExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
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

			"config": {
				Description:  "The information (in JSON format) managed by Terraform plan and apply.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newJson, _ := structure.NormalizeJsonString(new)
					oldJson, _ := structure.NormalizeJsonString(old)
					return newJson == oldJson
				},
			},

			"computed": {
				Description:  "A list of fields (in JSON format) returned from the executor script which are computed dynamically.",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsJSON,
			},

			"dynamic": {
				Description: "Fields (in JSON format) returned from the executor script which are computed dynamically.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"id_key": {
				Description:  "The name of the key which holds the unique identifier of the resource. e.g. 'id'",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"javascript": {
				Description:   "JavaScript to be executed internally by the Otto JavaScript interpreter.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"script", "executor"},
			},
		},
	}
}

func onCreate(d *schema.ResourceData, m interface{}) error {
	_, err := callExecutor("create", d, m)
	return err
}

func onRead(d *schema.ResourceData, m interface{}) error {
	_, err := callExecutor("read", d, m)
	return err
}

func onUpdate(d *schema.ResourceData, m interface{}) error {
	_, err := callExecutor("update", d, m)
	return err
}

func onDelete(d *schema.ResourceData, m interface{}) error {
	_, err := callExecutor("delete", d, m)
	return err
}

func onExists(d *schema.ResourceData, m interface{}) (bool, error) {
	return callExecutor("exists", d, m)
}

func getFromDefaultsOrResource(name string, defaults map[string]interface{}, dv interface{}, dok bool, required bool) (string, error) {
	//
	log.Printf("getFromDefaultsOrResource() field %s in %#v or %#v, %#v %#v\n", name, defaults, dv, dok, required)

	var result string
	found := false
	value, ok := defaults[name]
	if ok {
		if str, ok := value.(string); ok {
			result = str
			found = true
		}
	}

	if dok {
		str, ok := dv.(string)
		if ok {
			result = str
			found = true
		}
	}
	if found != true && required {
		return "", fmt.Errorf("missing required field %s in %v or %#v", name, defaults, dv)
	}
	return result, nil
}

// pickle Save some struct to a file for later unpickling
func pickle(event string, data interface{}) {

	// Open a file and dump JSON to it!
	fd, err := os.Create(event + ".json")
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(fd)
	err = enc.Encode(data)
	if err != nil {
		panic(err)
	}
	defer func() { _ = fd.Close() }()
}

// callExecutor - function to handle all the CRUDE. Returns with bool for 'exit'  all other responses
// are made in updates of the schema.ResourceData.
func callExecutor(event string, d ResourceLike, providerConfig interface{}) (bool, error) {

	effectiveDefaults, id, err := extractEssentialFields(event, d, providerConfig)
	if err != nil {
		return false, err
	}
	log.Printf("effectiveDefaults = %#v", effectiveDefaults)

	log.Printf("Executing: %s", event)
	configData, err := getConfigFromTF(d)
	if err != nil {
		return false, err
	}
	log.Printf("Executing: %s", string(configData))

	if _, ok := effectiveDefaults["javascript"]; ok {
		config := map[string]interface{}{}
		err := json.Unmarshal(configData, &config)
		if err != nil {
			return false, fmt.Errorf("expected JSON in 'config' but got: %#v", string(configData))
		}
		return callJavaScriptInterpreter(event, d, effectiveDefaults, id, config)
	}

	cmd := exec.Command(effectiveDefaults["executor"].(string), effectiveDefaults["script"].(string), event)
	cmd.Env = makeEnvironment(id, effectiveDefaults)

	if event == "delete" {
		cmd.Stdin = bytes.NewReader([]byte{})
	} else {
		cmd.Stdin = bytes.NewReader(configData)
	}
	// Call the executor
	rawResponse, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("Command error: %s\n", string(err.(*exec.ExitError).Stderr))
	}
	response, err := jsonSafeUnmarshal(rawResponse, err)
	if err != nil {
		return false, err
	}
	// Process the response
	if event == "exists" {
		var exists bool
		err = json.Unmarshal(rawResponse, &exists) // Need special unmarshall for atomic types
		if err != nil {
			return false, fmt.Errorf("expecting boolean from subprocess, got '%#v'", string(rawResponse))
		}
		return exists, nil
	} else if event == "delete" {
		d.SetId("")
	} else {
		responseMap, ok := response.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("expecting map[string]interface{} from subprocess, got '%#v'", string(rawResponse))
		}
		// Now move the computed fields into a separate "computed" attribute to avoid scrutiny by TF
		computedAsJSONbytes, err := moveComputedFields(effectiveDefaults, d, responseMap)
		if err != nil {
			return false, err
		}
		_ = d.Set("dynamic", string(computedAsJSONbytes))
		log.Printf("Executed: computed fields JSON is: %s", string(computedAsJSONbytes))

		// Get the id_key field from the response and move it into the special id member in the resourceData
		idKey := effectiveDefaults["id_key"].(string)
		if event == "create" {
			idRaw, ok := responseMap[idKey]
			if !ok {
				return false, fmt.Errorf("missing id attribute '%s' in response: %s", idKey, string(rawResponse))
			}
			id, ok := idRaw.(string)
			if !ok {
				return false, fmt.Errorf("expected string in id attribute '%s' in response but got: %#v", idKey, idRaw)
			}
			d.SetId(id)
		}
		delete(responseMap, idKey)

		// Now set the payload in the resource data 'config' field
		payloadBytes, err := json.Marshal(responseMap)
		if err != nil {
			return false, err
		}
		err = d.Set("config", string(payloadBytes))
		if err != nil {
			return false, err
		}
		log.Printf("Executed: setting data to: %s", string(payloadBytes))
	}

	return false, err
}

// moveComputedFields - Move the fields specified in effectiveDefaults["computed"] from responseMap into
// a returned map[]string encoded into JSON. Delete computed fields from the responseMap.
func moveComputedFields(effectiveDefaults map[string]interface{}, d ResourceLike, responseMap map[string]interface{}) ([]byte, error) {
	computed := make(map[string]interface{})
	cf := effectiveDefaults["computed"]
	log.Printf("Executed: computed fields: %s", cf)
	computedFields := make([]string, 3)
	cfjson, ok := cf.(string)
	if !ok {
		return nil, fmt.Errorf("expected string in 'computed' got %#v", cf)
	}
	err := json.Unmarshal([]byte(cfjson), &computedFields)
	if err != nil {
		return nil, fmt.Errorf("unable to parse JSON in 'computed' got %s", cfjson)
	}
	log.Printf("Executed: computed fields: %v", computedFields)

	dynamics := make(map[string]interface{})
	dynamic, ok := d.GetOk("dynamic")
	if ok {
		dynjson, ok := dynamic.(string)
		if ok {
			err = json.Unmarshal([]byte(dynjson), &dynamics)
			if err != nil {
				return nil, err
			}
		}
	}
	for _, name := range computedFields {
		cv, ok := responseMap[name]
		if ok {
			computed[name] = cv
			delete(responseMap, name)
		} else {
			computed[name] = dynamics[name]
		}
	}
	computedAsJSONbytes, err := json.Marshal(computed)
	return computedAsJSONbytes, err
}

// jsonSafeUnmarshal - copes with empty input
func jsonSafeUnmarshal(result []byte, err error) (interface{}, error) {
	var resource interface{}
	if len(result) == 0 {
		resource = nil
	} else {
		err = json.Unmarshal(result, &resource)
		if err != nil {
			return nil, err
		}
	}
	return resource, err
}

// makeEnvironment - Add the id and the 'environment' to the parent process environment,
// returning []string
func makeEnvironment(id string, effectiveDefaults map[string]interface{}) []string {
	environ := os.Environ()
	environ = append(environ, fmt.Sprintf("%s=%s", effectiveDefaults["id_key"], id))
	for k, v := range effectiveDefaults {
		if s, ok := v.(string); ok {
			e := fmt.Sprintf("%s=%s", k, s)
			environ = append(environ, e)
			log.Printf("Executing: with env var from default: %s", e)
		}
		if k == "environment" {
			if env, ok := v.(map[string]interface{}); ok {
				for envname, enval := range env {
					e := fmt.Sprintf("%s=%s", envname, enval)
					environ = append(environ, e)
					log.Printf("Executing: with env var from environment': %s", e)
				}
			}
		}
	}
	return environ
}

// extractEssentialFields - get the important fields from the provider config or resourceData.
// returning the a map[string] of the fields and the id field
func extractEssentialFields(event string, d ResourceLike, providerConfig interface{}) (map[string]interface{}, string, error) {
	essentialFields := map[string]bool{
		// map[field name]mandatory?
		"computed":    false,
		"environment": false,
		"executor":    false,
		"id_key":      true,
		"script":      false,
		"javascript":  false,
	}
	stringFields := []string{"id_key", "executor", "script"}

	var effectiveDefaults = map[string]interface{}{}

	id := d.Id()
	log.Printf("callExecutor() '%s' %s %#v", id, event, providerConfig)
	for n := range essentialFields {
		log.Printf("callExecutor() ResourceData field %s = %#v", n, d.Get(n))
	}
	// Validate provider configuration
	if providerConfig != nil {
		var ok bool
		effectiveDefaults, ok = providerConfig.(map[string]interface{})
		if !ok {
			return nil, "", fmt.Errorf("was expecting map[string]interface{} in provider configuration, got %#v", providerConfig)
		}
	}
	// Extract essential fields from provider configuration or resource data
	for k, required := range essentialFields {
		dv, dok := d.GetOk(k)
		value, err := getFromDefaultsOrResource(k, effectiveDefaults, dv, dok, required)
		if err != nil {
			return nil, "", err
		}
		effectiveDefaults[k] = value
		log.Printf("getFromDefaultsOrResource => field %s = %#v", k, value)
	}
	for _, stringFieldName := range stringFields {
		if _, ok := effectiveDefaults[stringFieldName].(string); !ok {
			return effectiveDefaults, id, fmt.Errorf("expected %s to be string, but got %#v", stringFieldName, effectiveDefaults[stringFieldName])
		}
	}
	return effectiveDefaults, id, nil
}

// getConfigFromTF - Validate and extract the 'config' JSON field from the resourceData, returning []byte
func getConfigFromTF(d ResourceLike) ([]byte, error) {
	dr, ok := d.GetOk("config")
	if !ok || dr == nil {
		return nil, fmt.Errorf("missing 'config'")
	}
	js, ok := dr.(string)
	if !ok {
		return nil, fmt.Errorf("expected string in 'config', but got: %#v", dr)
	}
	db := []byte(js)
	attributes := map[string]interface{}{}
	err := json.Unmarshal(db, &attributes)
	if err != nil {
		return nil, fmt.Errorf("expected JSON in 'config' but got: %#v", js)
	}
	configData, err := json.Marshal(attributes)
	return configData, err
}
