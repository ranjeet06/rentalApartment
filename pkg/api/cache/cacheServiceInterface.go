package cache

type ServiceCache interface {
	Set(key string, value interface{}) error
	Get(key string) (value interface{})
	Delete(key string) error
}
