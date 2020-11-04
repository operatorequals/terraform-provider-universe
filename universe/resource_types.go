package universe

// ResourceLike - An interface with the schema.ResourceData methods actually used in the provider
type ResourceLike interface {
	Id() string
	SetId(v string)
	Set(key string, value interface{}) error
	GetOk(key string) (interface{}, bool)
	Get(key string) interface{}
}
