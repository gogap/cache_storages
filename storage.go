package cache_storages

type StorageValue struct {
	V interface{} `json:"v"`
}

type CacheStorage interface {
	StorageType() string

	SetObject(Key string, v interface{}) (err error)
	GetObject(key string, v interface{}) (err error)
	GetMultiObject(keyValues map[string]interface{}) (err error)

	Set(key string, v string) (err error)
	Get(key string) (v string, err error)
	GetMulti(key string) (values map[string]string, err error)

	Touch(key string, seconds int32) (err error)

	Increment(key string, delta uint64) (newValue uint64, err error)
	Decrement(key string, delta uint64) (newValue uint64, err error)

	Delete(key string) error
	DeleteAll() error
}
