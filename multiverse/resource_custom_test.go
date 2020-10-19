package multiverse

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"os"
	"testing"
)

func Test_getFromDefaultsOrResource(t *testing.T) {
	_, err := getFromDefaultsOrResource("foo", map[string]interface{}{"foo": "12"}, &schema.ResourceData{})
	if err == nil {
		t.Fail()
	}
	_, err = getFromDefaultsOrResource("foo", map[string]interface{}{"environment": map[string]interface{}{"foo": "12"}}, &schema.ResourceData{})
	if err != nil {
		t.Fail()
	}
	rd := schema.ResourceData{}
	// NOW WHAT
	rd.Set("foo", "44")
	_, err = getFromDefaultsOrResource("foo", map[string]interface{}{"foo": "12"}, &rd)
	//if err != nil {
	//	t.Fail()
	//}
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
