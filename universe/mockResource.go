package universe

// A struct that conforms to the ResourceLike interface used for testing.
type mockResource struct {
	id     string
	fields map[string]interface{}
}

func (d *mockResource) Id() string {
	return d.id
}
func (d *mockResource) SetId(v string) {
	d.id = v
}
func (d *mockResource) Set(key string, value interface{}) error {
	d.fields[key] = value
	return nil
}
func (d *mockResource) GetOk(key string) (interface{}, bool) {
	v, ok := d.fields[key]
	return v, ok
}
func (d *mockResource) Get(key string) interface{} {
	v, ok := d.GetOk(key)
	if !ok {
		v = ""
	}
	return v
}

// NewMockResource - Make Fake
func NewMockResource() ResourceLike {
	var d mockResource
	d.fields = map[string]interface{}{}
	return &d
}
