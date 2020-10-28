package multiverse

import (
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
	err := testAccProvider.Configure(terraform.NewResourceConfigRaw(map[string]interface{}{}))
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
