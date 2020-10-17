package multiverse

import (
	"testing"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// var testProviders = map[string]schema.Resource{
// 	"multiverse": Provider(),
// }
