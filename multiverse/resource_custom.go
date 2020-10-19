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

			"data": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newJson, _ := structure.NormalizeJsonString(new)
					oldJson, _ := structure.NormalizeJsonString(old)
					return newJson == oldJson
				},
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
	return do("create", d, m)
}

func onRead(d *schema.ResourceData, m interface{}) error {
	return do("read", d, m)
}

func onUpdate(d *schema.ResourceData, m interface{}) error {
	return do("update", d, m)
}

func onDelete(d *schema.ResourceData, m interface{}) error {
	return do("delete", d, m)
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
	if found != true {
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
	defer fd.Close()
}
func do(event string, d *schema.ResourceData, defaults interface{}) error {
	//
	id := d.Id()
	for _, n := range []string{"script", "executor", "id_key"} {
		log.Printf("do() ResourceData field %s = %#v", n, d.Get(n))
	}

	def, ok := defaults.(map[string]interface{})
	if !ok {
		return fmt.Errorf("was expecting map[string]interface{} in 'environment', got %v", defaults)
	}
	for _, k := range []string{"id_key", "executor", "script"} {
		value, err := getFromDefaultsOrResource(k, def, d)
		if err != nil {
			return err
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
	dr, ok := d.GetOk("data")
	if !ok || dr == nil {
		return fmt.Errorf("bad JSON in script: %v", dr)
	}
	js, ok := dr.(string)
	if !ok {
		return fmt.Errorf("bad JSON in script: %v", dr)
	}
	db := []byte(js)
	attributes := map[string]interface{}{}
	err := json.Unmarshal(db, &attributes)
	if err != nil {
		return fmt.Errorf("bad JSON in script: %s", js)
	}
	datab, err := json.Marshal(attributes)
	if err != nil {
		return err
	}
	log.Printf("Executing: %s", string(datab))

	cmd := exec.Command(def["executor"].(string), def["script"].(string), event)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", idKey, id))
	for k, v := range def {
		if s, ok := v.(string); ok {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, s))
		}
		if env, ok := v.(map[string]interface{}); ok {
			for envname, enval := range env {
				cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", envname, enval))
			}
		}
	}
	for _, e := range cmd.Env {
		log.Printf("Executing: environ: %s", e)
	}

	if event == "delete" {
		cmd.Stdin = bytes.NewReader([]byte(id))
	} else {
		cmd.Stdin = bytes.NewReader(datab)
	}

	result, err := cmd.Output()

	if err != nil {
		log.Printf("Command error: %s\n", string(err.(*exec.ExitError).Stderr))
		return err
	}

	var resource interface{}
	err = json.Unmarshal(result, &resource)
	if err != nil {
		return err
	}
	if event == "delete" {
		d.SetId("")
	} else {
		rm, ok := resource.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expecting map[string]interface{} from subprocess, got '%#v'", resource)
		}
		_, ok = rm[idKey]
		if !ok && event == "create" {
			return fmt.Errorf("missing id attribute '%s' in response: %s", idKey, result)
		}
		if event == "create" {
			d.SetId(rm[idKey].(string))
		}
		delete(rm, idKey)
		resultb, err := json.Marshal(rm)
		if err != nil {
			return err
		}
		err = d.Set("data", string(resultb))
	}

	return err
}
