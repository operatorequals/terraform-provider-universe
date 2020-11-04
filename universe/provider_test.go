package universe

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderConfigure(t *testing.T) {
	d := NewMockResource()
	_ = d.Set("script", "S")
	_ = d.Set("environment", map[string]interface{}{
		"X": "2",
	})
	p, err := providerConfigure(d)
	if err != nil {
		t.FailNow()
	}
	pmap, ok := p.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	script, ok := pmap["script"].(string)
	if !ok {
		t.FailNow()
	}
	if script != "S" {
		t.Fail()
	}
	env, ok := pmap["environment"]
	if !ok {
		t.FailNow()
	}
	envmap, ok := env.(map[string]interface{})
	x, ok := envmap["X"]
	if !ok {
		t.Fail()
	}
	xs, ok := x.(string)
	if !ok {
		t.Fail()
	}
	if xs != "2" {
		t.Fail()
	}
}

func TestProviderNameFromBinaryOrEnvironment(t *testing.T) {
	pn := "marigolds"
	err := os.Setenv(EnvProviderNameVar, pn)
	if err != nil {
		t.Fail()
	}
	n := getProviderNameFromBinaryOrEnvironment()
	if n != pn {
		t.Fail()
	}
	err = os.Unsetenv(EnvProviderNameVar)
	if err != nil {
		t.Fail()
	}
	t.Logf("binary name is %s", filepath.Base(os.Args[0]))
	n = getProviderNameFromBinaryOrEnvironment()
	if n != DefaultProviderName {
		t.Fail()
	}
}

func TestGetResourceTypeNamesFromEnvironment(t *testing.T) {
	providerName := "dugong"
	resourceNamesName := "TERRAFORM_" + strings.ToUpper(providerName) + "_RESOURCETYPES"
	rn := "cats dogs     dugong_hats"
	err := os.Setenv(resourceNamesName, rn)
	if err != nil {
		t.Fail()
	}
	names := getResourceTypeNamesFromEnvironment(providerName)
	t.Logf("%v", names)
	if !reflect.DeepEqual(names, map[string]bool{"dugong_cats": true, "dugong_dogs": true, "dugong_hats": true, "dugong": true}) {
		t.Fail()
	}
}

func TestAccPreCheck(t *testing.T) {

	testAccProvider := Provider()
	err := testAccProvider.Configure(nil, terraform.NewResourceConfigRaw(map[string]interface{}{}))
	if err != nil {
		t.Fatal(err)
	}
}
func TestDiffSuppressComputed(t *testing.T) {

	if false == diffSuppressComputed("k", "", "", nil) {
		t.Fail()
	}

	if true == diffSuppressComputed("k", `{"@A": "23", "B":44 }`, `{"@A": 12}`, nil) {
		t.Fail()
	}

	if true == diffSuppressComputed("k", `{"@A": "23", "B":44 }`, `{"@A": 12, "B": "44" }`, nil) {
		t.Fail()
	}

}
