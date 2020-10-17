package multiverse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
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
			"executor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"script": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"data": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"id_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"resource": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
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

func do(event string, d *schema.ResourceData, m interface{}) error {
	//
	log.Printf("Executing: %s %s %s %#v", d.Get("executor"), d.Get("script"), event, d.Get("data"))

	cmd := exec.Command(d.Get("executor").(string), d.Get("script").(string), event)

	if event == "delete" {
		cmd.Stdin = bytes.NewReader([]byte(d.Id()))
	} else {
		d := []byte(d.Get("data").(string))
		var ignore interface{}
		err := json.Unmarshal(d, &ignore)
		if err != nil {
			// User gave bad JSON
			return fmt.Errorf("bad JSON in script: %s", string(d))
		}
		cmd.Stdin = bytes.NewReader(d)
	}

	result, err := cmd.Output()

	if err != nil {
		log.Printf("Command error: %s\n", string(err.(*exec.ExitError).Stderr))
		return err
	}

	var resource map[string]interface{}
	err = json.Unmarshal([]byte(result), &resource)
	if err != nil {
		return err
	}
	if event == "delete" {
		d.SetId("")
	} else {
		key := d.Get("id_key").(string)
		d.Set("resource", resource)
		d.SetId(resource[key].(string))
	}

	return err
}
