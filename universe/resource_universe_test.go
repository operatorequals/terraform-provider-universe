package universe

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
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
		"executor": "python3",
		"script":   "resource_universe_test.py",
	}
	_, err := callExecutor("create", d, config)
	if err != nil {
		t.FailNow()
	}
	if d.Id() != "42" {
		t.Fail()
	}
	c := d.Get("config")
	n1, _ := structure.NormalizeJsonString(c)
	n2, _ := structure.NormalizeJsonString(`{"@created":"26/10/2020 18:55:51", "album":"white"}`)
	if n1 != n2 {
		t.Fail()
	}
}
func Test_callExecutorUpdate(t *testing.T) {
	d := NewMockResource()
	d.SetId("42")
	_ = d.Set("config", `{"album": "black"}`)
	config := map[string]interface{}{
		"id_key":   "id",
		"executor": "python3",
		"script":   "resource_universe_test.py",
	}
	_, err := callExecutor("update", d, config)
	if err != nil {
		t.FailNow()
	}
	c := d.Get("config")
	n1, _ := structure.NormalizeJsonString(c)
	n2, _ := structure.NormalizeJsonString(`{"@created":"26/10/2020 18:55:51", "album":"black"}`)
	if n1 != n2 {
		t.Fail()
	}
}

func Test_callExecutorExists(t *testing.T) {
	d := NewMockResource()
	_ = d.Set("config", `{"album": "white"}`)
	config := map[string]interface{}{
		"id_key":   "id",
		"executor": "python3",
		"script":   "resource_universe_test.py",
	}
	d.SetId("42")
	exists, err := callExecutor("exists", d, config)
	if !exists || err != nil {
		t.Fail()
	}
}

func Test_callExecutorDelete(t *testing.T) {
	d := NewMockResource()
	d.SetId("42")
	_ = d.Set("config", `{"album": "white"}`)
	config := map[string]interface{}{
		"id_key":   "id",
		"executor": "python3",
		"script":   "resource_universe_test.py",
	}
	_, err := callExecutor("delete", d, config)
	if err != nil {
		t.FailNow()
	}
}

func Test_callExecutorBad(t *testing.T) {
	d := NewMockResource()
	_ = d.Set("config", `{"album": "white"}`)
	config := map[string]interface{}{
		"id_key":   "id",
		"executor": "", // Bad or wrong path to program
		"script":   "resource_universe_test.py",
	}
	_, err := callExecutor("create", d, config)
	if err == nil {
		t.FailNow()
	}
}
