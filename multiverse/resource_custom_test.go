package multiverse

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test_getFromDefaultsOrResource(t *testing.T) {
	value, err := getFromDefaultsOrResource("foo", map[string]interface{}{"foo": "12"}, "", false, true)
	if err != nil || value != "12" {
		t.Fail()
	}
	value, err = getFromDefaultsOrResource("foo", map[string]interface{}{"foo": "12"}, "44", true, true)
	if err != nil || value != "44" {
		t.Fail()
	}
	value, err = getFromDefaultsOrResource("foo", map[string]interface{}{"bar": "12"}, "44", false, true)
	if err == nil {
		t.Fail()
	}
	value, err = getFromDefaultsOrResource("foo", map[string]interface{}{"bar": "12"}, "44", false, false)
	if err != nil {
		t.Fail()
	}
}

func Test_pickle(t *testing.T) {
	data :=
		map[string]interface{}{
			"alpha":   "A",
			"bravo":   "2",
			"charlie": true,
		}
	tmpfile, err := ioutil.TempFile("", "test-pickle-")
	if err != nil {
		t.FailNow()
	}
	pickle(tmpfile.Name(), data)

	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name()) // clean up

}
