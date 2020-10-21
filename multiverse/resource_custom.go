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
				Type:     schema.TypeString,
				Optional: true,
			},

			"script": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"config": {
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"id_key": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func onCreate(d *schema.ResourceData, m interface{}) error {
	_, err := do("create", d, m)
	return err
}

func onRead(d *schema.ResourceData, m interface{}) error {
	_, err := do("read", d, m)
	return err
}

func onUpdate(d *schema.ResourceData, m interface{}) error {
	_, err := do("update", d, m)
	return err
}

func onDelete(d *schema.ResourceData, m interface{}) error {
	_, err := do("delete", d, m)
	return err
}

func onExists(d *schema.ResourceData, m interface{}) (bool, error) {
	return do("exists", d, m)
}

func getFromDefaultsOrResource(name string, defaults map[string]interface{}, d *schema.ResourceData) (string, error) {
	//
	log.Printf("getFromDefaultsOrResource() field %s in %#v or %#v", name, defaults, d)

	var result string
	found := false
	value, ok := defaults[name]
	if ok {
		if str, ok := value.(string); ok {
			result = str
			found = true
		}
		if m, ok := value.(map[string]interface{}); ok {
			if value, ok = m[name]; ok {
				if str, ok := value.(string); ok {
					result = str
					found = true
				}
			}
		}
	}
	x, ok := d.GetOk(name)
	if ok {
		log.Printf("getFromDefaultsOrResource(98) field %s = %#v", name, x)
		str, ok := x.(string)
		if ok {
			log.Printf("getFromDefaultsOrResource(101) field %s = %#v", name, str)
			result = str
			found = true
		}
		log.Printf("getFromDefaultsOrResource(105) field %s = %#v", name, str)
	}
	if found != true && name != "computed" { // TODO
		return "", fmt.Errorf("missing required field %s in %v or %v", name, defaults, x)
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
func do(event string, d *schema.ResourceData, defaults interface{}) (bool, error) {
	// TODO Make nicer code
	id := d.Id()
	for _, n := range []string{"script", "executor", "id_key"} {
		log.Printf("do() ResourceData field %s = %#v", n, d.Get(n))
	}

	def, ok := defaults.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("was expecting map[string]interface{} in 'environment', got %v", defaults)
	}
	for _, k := range []string{"id_key", "executor", "script", "computed"} {
		value, err := getFromDefaultsOrResource(k, def, d)
		if err != nil {
			return false, err
		}
		def[k] = value
		log.Printf("getFromDefaultsOrResource => field %s = %#v", k, value)
	}
	var idKey string
	if idk, ok := def["id_key"]; ok {
		idKey = idk.(string)
	} else {
		idKey = "no id key!"
	}
	log.Printf("Executing: %s", event)
	dr, ok := d.GetOk("config")
	if !ok || dr == nil {
		return false, fmt.Errorf("bad JSON in script: %v", dr)
	}
	js, ok := dr.(string)
	if !ok {
		return false, fmt.Errorf("bad JSON in script: %v", dr)
	}
	db := []byte(js)
	attributes := map[string]interface{}{}
	err := json.Unmarshal(db, &attributes)
	if err != nil {
		return false, fmt.Errorf("bad JSON in script: %s", js)
	}
	datab, err := json.Marshal(attributes)
	if err != nil {
		return false, err
	}
	log.Printf("Executing: %s", string(datab))

	cmd := exec.Command(def["executor"].(string), def["script"].(string), event)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", idKey, id))
	for k, v := range def {
		if s, ok := v.(string); ok {
			e := fmt.Sprintf("%s=%s", k, s)
			cmd.Env = append(cmd.Env, e)
			log.Printf("Executing: environ default: %s", e)
		}
		if k == "environment" {
			if env, ok := v.(map[string]interface{}); ok {
				for envname, enval := range env {
					e := fmt.Sprintf("%s=%s", envname, enval)
					cmd.Env = append(cmd.Env, e)
					log.Printf("Executing: environ environment: %s", e)
				}
			}
		}
	}

	if event == "delete" {
		cmd.Stdin = bytes.NewReader([]byte{})
	} else {
		cmd.Stdin = bytes.NewReader(datab)
	}

	result, err := cmd.Output()

	if err != nil {
		log.Printf("Command error: %s\n", string(err.(*exec.ExitError).Stderr))
		return false, err
	}

	var resource interface{}
	if len(result) == 0 {
		resource = nil
	} else {
		err = json.Unmarshal(result, &resource)
		if err != nil {
			return false, err
		}
	}
	if event == "exists" {
		var exists bool
		err = json.Unmarshal(result, &exists) // Need special unmarshall for atomic types
		if err != nil {
			return false, fmt.Errorf("expecting boolean from subprocess, got '%#v'", string(result))
		}
		return exists, nil
	} else if event == "delete" {
		d.SetId("")
	} else {
		rm, ok := resource.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("expecting map[string]interface{} from subprocess, got '%#v'", string(result))
		}
		// Now move the computed fields into a separate "computed" attribute to avoid scrutiny by TF
		computed := make(map[string]interface{})
		cf := def["computed"]
		log.Printf("Executed: computed fields: %s", cf)
		computedFields := make([]string, 3)
		cfjson, ok := cf.(string)
		if !ok {
			return false, fmt.Errorf("unable to get string in 'computed' got %v", cf)
		}
		err := json.Unmarshal([]byte(cfjson), &computedFields)
		if err != nil {
			return false, fmt.Errorf("unable to parse 'computed' got %s", cfjson)
		}
		dynamics := make(map[string]interface{})
		dynamic, ok := d.GetOk("dynamic")
		if ok {
			dynjson, ok := dynamic.(string)
			if ok {
				err = json.Unmarshal([]byte(dynjson), &dynamics)
				if err != nil {
					return false, err
				}
			}
		}
		log.Printf("Executed: computed fields: %v", computedFields)
		for _, name := range computedFields {
			cv, ok := rm[name]
			if ok {
				computed[name] = cv
				delete(rm, name)
			} else {
				computed[name] = dynamics[name]
			}
		}
		computedAsJSONbytes, err := json.Marshal(computed)
		_ = d.Set("dynamic", string(computedAsJSONbytes))
		log.Printf("Executed: computed fields JSON is: %s", string(computedAsJSONbytes))

		_, ok = rm[idKey]
		if !ok && event == "create" {
			return false, fmt.Errorf("missing id attribute '%s' in response: %s", idKey, string(result))
		}
		if event == "create" {
			d.SetId(rm[idKey].(string))
		}
		delete(rm, idKey)
		resultb, err := json.Marshal(rm)
		if err != nil {
			return false, err
		}
		err = d.Set("config", string(resultb))
		log.Printf("Executed: setting data to: %s", string(resultb))
	}

	return false, err
}
