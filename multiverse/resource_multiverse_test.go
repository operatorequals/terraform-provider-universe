package multiverse

import (
	"fmt"
	"reflect"
	"testing"
)

// (name string, defaults map[string]interface{}, d ResourceLike, required bool

func Test_getFromDefaultsOrResource(t *testing.T) {
	r := NewMockResource()
	value, ok := getFromDefaultsOrResource("foo", map[string]interface{}{"foo": "12"}, r, true)
	if !ok || value != "12" {
		t.Fail()
	}
	r = &mockResource{id: "", fields: map[string]interface{}{"foo": "44"}}
	value, ok = getFromDefaultsOrResource("foo", map[string]interface{}{"foo": "12"}, r, true)
	if !ok || value != "44" {
		t.Fail()
	}
	r = NewMockResource()
	_ = r.Set("bar", 44)
	value, ok = getFromDefaultsOrResource("foo", map[string]interface{}{"bar": "12"}, r, true)
	if ok {
		t.Fail()
	}
}

func Test_callExecutorCreate(t *testing.T) {
	d := NewMockResource()
	_ = d.Set("config", `{"album": "white"}`)
	config := map[string]interface{}{
		"id_key":   "id",
		"computed": `["created"]`,
		"executor": "python3",
		"script":   "resource_multiverse_test.py",
	}
	_, err := callExecutor("create", d, config)
	if err != nil {
		t.FailNow()
	}
	if d.id != "42" {
		t.Fail()
	}
	c := d.Get("config")
	if !reflect.DeepEqual(c, `{"album":"white"}`) {
		t.Fail()
	}
	dyn := d.Get("dynamic")
	if !reflect.DeepEqual(dyn, `{"created":"26/10/2020 18:55:51"}`) {
		t.Fail()
	}
}
func Test_callExecutorUpdate(t *testing.T) {
	d := NewMockResource()
	d.SetId("42")
	_ = d.Set("config", `{"album": "black"}`)
	config := map[string]interface{}{
		"id_key":   "id",
		"computed": `["created"]`,
		"executor": "python3",
		"script":   "resource_multiverse_test.py",
	}
	_, err := callExecutor("update", d, config)
	if err != nil {
		t.FailNow()
	}
	c := d.Get("config")
	if !reflect.DeepEqual(c, `{"album":"black"}`) {
		t.Fail()
	}
	dyn := d.Get("dynamic")
	if !reflect.DeepEqual(dyn, `{"created":"26/10/2020 18:55:51"}`) {
		t.Fail()
	}
}

func Test_callExecutorExists(t *testing.T) {
	d := NewMockResource()
	_ = d.Set("config", `{"album": "white"}`)
	config := map[string]interface{}{
		"id_key":   "id",
		"computed": `["created"]`,
		"executor": "python3",
		"script":   "resource_multiverse_test.py",
	}
	d.SetId("42")
	exists, err := callExecutor("exists", d, config)
	fmt.Printf("%#v %#v", exists, err)
	if !exists {
		t.Fail()
	}
}

func Test_callExecutorDelete(t *testing.T) {
	d := NewMockResource()
	d.SetId("42")
	_ = d.Set("config", `{"album": "white"}`)
	config := map[string]interface{}{
		"id_key":   "id",
		"computed": `["created"]`,
		"executor": "python3",
		"script":   "resource_multiverse_test.py",
	}
	_, err := callExecutor("delete", d, config)
	if err != nil {
		t.FailNow()
	}
}
