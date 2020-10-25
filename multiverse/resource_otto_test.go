package multiverse

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_callJavaScriptInterpreter(t *testing.T) {
	d := NewMockResource()
	config := map[string]interface{}{
		"album": "white",
	}
	effectiveDefaults := map[string]interface{}{
		"id_key":     "id",
		"javascript": "../examples/javascript/echo.js",
		"computed":   `["created"]`,
	}
	id := "1"
	v, err := callJavaScriptInterpreter("create", d, effectiveDefaults, id, config)
	fmt.Printf("%#v %#v", v, err)
	if d.id != "42" {
		t.Fail()
	}
	c := d.Get("config")
	if !reflect.DeepEqual(c, `{"album":"white"}`) {
		t.Fail()
	}
}
