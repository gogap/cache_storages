package cache_storages

import (
	"encoding/json"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedStorage struct {
	client *memcache.Client
}

func NewMemcachedStorage(endpoints ...string) (storage *MemcachedStorage, err error) {
	client := memcache.New(endpoints...)
	if client == nil {
		err = fmt.Errorf("create storage client failed, %v", endpoints)
		return
	}

	storage = new(MemcachedStorage)
	storage.client = client
	return
}

func (p *MemcachedStorage) StorageType() string {
	return "memcached"
}

func (p *MemcachedStorage) SetObject(key string, v interface{}) (err error) {
	sV := StorageValue{V: v}
	if bJsonV, e := json.Marshal(sV); e != nil {
		err = fmt.Errorf("marshal object to json failed, key: %s, value: %v, %s", key, v, e)
		return
	} else {
		item := &memcache.Item{Key: key, Value: bJsonV}
		if e := p.client.Set(item); e != nil {
			err = fmt.Errorf("set key-value to memcache failed, key:%s, value: %v, %s", key, v, e)
			return
		}
	}
	return
}

func (p *MemcachedStorage) GetObject(key string, v interface{}) (err error) {
	var item *memcache.Item
	if item, err = p.client.Get(key); err != nil {
		err = fmt.Errorf("get key failed, key: %s, %s", key, err)
		return
	}

	bJsonV := item.Value
	if bJsonV == nil {
		v = nil
		return
	}

	sv := StorageValue{V: v}
	if e := json.Unmarshal(bJsonV, &sv); e != nil {
		err = fmt.Errorf("unmarshal json to object failed, key: %s, value: %v, %s", key, bJsonV, e)
		return
	}

	return
}

func (p *MemcachedStorage) Set(key string, v string) (err error) {
	item := &memcache.Item{Key: key, Value: []byte(v)}
	if e := p.client.Set(item); e != nil {
		err = fmt.Errorf("set key-value to memcache failed, key:%s, value: %s, %s", key, v, e)
		return
	}
	return
}
func (p *MemcachedStorage) Get(key string) (v string, err error) {
	var item *memcache.Item
	if item, err = p.client.Get(key); err != nil {
		err = fmt.Errorf("get key failed, key: %s, %s", key, err)
		return
	} else if item.Value != nil {
		v = string(item.Value)
		return
	}
	v = ""
	return
}

func (p *MemcachedStorage) Touch(key string, seconds int32) (err error) {
	if err = p.client.Touch(key, seconds); err != nil {
		err = fmt.Errorf("touch key failed, key: %s, seconds: %d, %s", key, seconds, err)
		return
	}
	return
}

func (p *MemcachedStorage) GetMulti(keys []string) (values map[string]string, err error) {
	if items, e := p.client.GetMulti(keys); e != nil {
		err = fmt.Errorf("get multi keys error, keys: %v, %s", keys, e)
		return
	} else {
		values = make(map[string]string, len(items))
		for k, v := range items {
			if v == nil {
				values[k] = ""
			} else {
				values[k] = string(v.Value)
			}
		}
	}
	return
}

func (p *MemcachedStorage) GetMultiObject(keyValues map[string]interface{}) (err error) {

	keys := []string{}
	for key, value := range keyValues {
		if value == nil {
			err = fmt.Errorf("values did not contain the interface{} of the key: %s", key)
			return
		}
		keys = append(keys, key)
	}

	if items, e := p.client.GetMulti(keys); e != nil {
		err = fmt.Errorf("get multi keys failed, keys: %v, %s", keys, err)
		return
	} else {
		for _, item := range items {
			bJsonV := item.Value
			if bJsonV == nil {
				keyValues[item.Key] = nil
				continue
			}
			v, _ := keyValues[item.Key]
			sv := StorageValue{V: v}
			if e := json.Unmarshal(bJsonV, &sv); e != nil {
				keyValues = nil
				err = fmt.Errorf("unmarshal json to object failed, key: %s, value: %v, %s", item.Key, bJsonV, e)
				return
			}
		}
	}
	return
}

func (p *MemcachedStorage) Increment(key string, delta uint64) (newValue uint64, err error) {
	if newValue, err = p.client.Increment(key, delta); err != nil {
		err = fmt.Errorf("increment key failed, key: %s, delta: %d, %s", key, delta, err)
		return
	}
	return
}
func (p *MemcachedStorage) Decrement(key string, delta uint64) (newValue uint64, err error) {
	if newValue, err = p.client.Decrement(key, delta); err != nil {
		err = fmt.Errorf("decrement key failed, key: %s, delta: %d, %s", key, delta, err)
		return
	}
	return
}

func (p *MemcachedStorage) Delete(key string) (err error) {
	if err = p.client.Delete(key); err != nil {
		err = fmt.Errorf("delete key failed, key: %s, %s", key, err)
		return
	}
	return
}
func (p *MemcachedStorage) DeleteAll() (err error) {
	if err = p.client.DeleteAll(); err != nil {
		err = fmt.Errorf("delete all failed, %s", err)
		return
	}
	return
}
