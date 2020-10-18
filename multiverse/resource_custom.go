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
				Required: true,
			},

			"script": {
				Type:     schema.TypeString,
				Required: true,
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
				Required:     true,
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

func do(event string, d *schema.ResourceData, config interface{}) error {
	//
	configuration, ok := config.(string)
	if !ok {
		return fmt.Errorf("bad configuration: %v", config)
	}
	id := d.Id()
	idKey := d.Get("id_key").(string)
	command := d.Get("executor").(string)
	script := d.Get("script").(string)

	log.Printf("Executing: %s %s %s %s", event, command, script, configuration)
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
	attributes[idKey] = id
	datab, err := json.Marshal(attributes)
	if err != nil {
		return err
	}
	log.Printf("Executing: %s", string(datab))

	cmd := exec.Command(command, script, event)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "multiverse="+configuration)

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

	var resource map[string]interface{}
	err = json.Unmarshal(result, &resource)
	if err != nil {
		return err
	}
	if event == "delete" {
		d.SetId("")
	} else {
		_, ok := resource[idKey]
		if !ok {
			return fmt.Errorf("missing id attribute '%s' in response: %s", idKey, result)
		}
		d.SetId(resource[idKey].(string))
		delete(resource, idKey)
		resultb, err := json.Marshal(resource)
		if err != nil {
			return err
		}
		err = d.Set("data", string(resultb))
	}

	return err
}
