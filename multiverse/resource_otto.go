package multiverse

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// callJavaScriptInterpreter - function to handle all the CRUDE in the embedded otto interpreter
func callJavaScriptInterpreter(event string, d ResourceLike, effectiveDefaults map[string]interface{}, id string, config map[string]interface{}) (bool, error) {
	vm := otto.New()
	idKey := effectiveDefaults["id_key"].(string)
	err := vm.Set(idKey, id)
	if err != nil {
		return false, err
	}
	for k, v := range effectiveDefaults {
		err := vm.Set(k, v)
		if err != nil {
			return false, err
		}
	}
	err = vm.Set("config", config)
	if err != nil {
		return false, err
	}
	err = vm.Set("event", event)
	if err != nil {
		return false, err
	}
	pwd, _ := os.Getwd()
	path, err := filepath.Abs(pwd + "/" + effectiveDefaults["javascript"].(string))
	if err != nil {
		return false, err
	}
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	responseVal, err := vm.Run(code)
	if err != nil {
		return false, err
	}
	response, _ := responseVal.Export()
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("expecting map[string]interface{} from subprocess, got '%#v'", response)
	}
	// Now move the computed fields into a separate "computed" attribute to avoid scrutiny by TF
	computedAsJSONbytes, err := moveComputedFields(effectiveDefaults, d, responseMap)
	if err != nil {
		return false, err
	}
	_ = d.Set("dynamic", string(computedAsJSONbytes))
	log.Printf("Executed: computed fields JSON is: %s", string(computedAsJSONbytes))

	// Get the id_key field from the response and move it into the special id member in the resourceData
	if event == "create" {
		idRaw, ok := responseMap[idKey]
		if !ok {
			return false, fmt.Errorf("missing id attribute '%s' in response: %#v", idKey, responseMap)
		}
		id, ok := idRaw.(string)
		if !ok {
			return false, fmt.Errorf("expected string in id attribute '%s' in response but got: %#v", idKey, idRaw)
		}
		d.SetId(id)
	}
	delete(responseMap, idKey)

	// Now set the payload in the resource data
	payloadBytes, err := json.Marshal(responseMap)
	if err != nil {
		return false, err
	}
	err = d.Set("config", string(payloadBytes))
	if err != nil {
		return false, err
	}
	log.Printf("Executed: set config to: %s", string(payloadBytes))

	return true, err
}
