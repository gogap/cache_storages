package cache_storages

import (
	"testing"
	//"time"
)

func TestRedisSetNX(t *testing.T) {
	storage, err := NewRedisStorage("127.0.0.1:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}

	result, err := storage.SetNX("key111sdfsdf", "value")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestRedisGetSetNX(t *testing.T) {
	storage, err := NewRedisStorage("127.0.0.1:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}

	result, err := storage.GetSet("key111sdfsdf", "value1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

/*
func TestRedisGetSetObject(t *testing.T) {
	type value struct {
		Name string
		Year int
	}

	var storage CacheStorage
	var err error
	storage, err = NewRedisStorage("127.0.0.1:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}

	var v value
	v.Name = "y"
	v.Year = 24

	err = storage.SetObject("key", v, 10)
	if err != nil {
		t.Error(err)
		return
	}

	var reply value
	err = storage.GetObject("key", &reply)
	if err != nil {
		t.Error(err)
		return
	}

	if reply.Name != v.Name || reply.Year != v.Year {
		t.Error("get object error", reply)
		return
	}
}

func TestRedisGetSet(t *testing.T) {
	storage, err := NewRedisStorage("127.0.0.1:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}

	err = storage.Set("key", "value", 1)
	if err != nil {
		t.Error(err)
		return
	}

	//验证set get
	value, err := storage.Get("key")
	if err != nil {
		t.Error(err)
		return
	}

	if value != "value" {
		t.Error("get string error", value)
		return
	}

	time.Sleep(time.Second * 3)
	//验证超时
	value, err = storage.Get("key")
	if err != nil {
		t.Error(err)
		return
	}

	if value != "" {
		t.Error("get string error", value)
		return
	}

	err = storage.Set("key", "value", 0)
	if err != nil {
		t.Error(err)
		return
	}

	value, err = storage.Get("key")
	if err != nil {
		t.Error(err)
		return
	}

	if value != "value" {
		t.Error("get string error", value)
		return
	}

	err = storage.Touch("key", 1)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second * 3)
	value, err = storage.Get("key")
	if err != nil {
		t.Error(err)
		return
	}

	if value != "" {
		t.Error("get string error", value)
		return
	}

}

func TestRedisGetSetMultiObject(t *testing.T) {
	type value struct {
		Name string
		Year int
	}

	storage, err := NewRedisStorage("127.0.0.1:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}

	var v value
	v.Name = "y"
	v.Year = 24

	err = storage.SetObject("key", v, 10)
	if err != nil {
		t.Error(err)
		return
	}

	var v2 value
	v2.Name = "l"
	v2.Year = 48

	err = storage.SetObject("key2", v2, 10)
	if err != nil {
		t.Error(err)
		return
	}

	var vs = make(map[string]interface{})
	vs["key"] = new(value)
	vs["key2"] = new(value)
	err = storage.GetMultiObject(vs)
	if err != nil {
		t.Error(err)
		return
	}

	if value, ok := vs["key"].(*value); ok {
		if value.Name != "y" || value.Year != 24 {
			t.Error("value error")
			return
		}
	} else {
		t.Error("type error")
		return
	}

	if value, ok := vs["key2"].(*value); ok {
		if value.Name != "l" || value.Year != 48 {
			t.Error("value error")
			return
		}
	} else {
		t.Error("type error")
		return
	}
}

func TestRedisGetSetMulti(t *testing.T) {
	storage, err := NewRedisStorage("127.0.0.1:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}

	err = storage.Set("key", "value", 10)
	if err != nil {
		t.Error(err)
		return
	}

	err = storage.Set("key2", "value2", 10)
	if err != nil {
		t.Error(err)
		return
	}

	values, err := storage.GetMulti([]string{"key", "key2"})
	if err != nil {
		t.Error(err)
		return
	}
	if values["key"] != "value" ||
		values["key2"] != "value2" {
		t.Error("value error")
		return
	}
}

func TestRedisDelete(t *testing.T) {
	storage, err := NewRedisStorage("127.0.0.1:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}
	err = storage.Set("key", "value", 10)
	if err != nil {
		t.Error(err)
		return
	}
	err = storage.Delete("key")
	if err != nil {
		t.Error(err)
		return
	}
	value, err := storage.Get("key")
	if err != nil {
		t.Error(err)
		return
	}
	if value != "" {
		t.Error("value error")
		return
	}
}

func TestRedisDeleteAll(t *testing.T) {
	storage, err := NewRedisStorage("127.0.0.1:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}
	err = storage.Set("key", "value", 10)
	if err != nil {
		t.Error(err)
		return
	}
	err = storage.Set("key2", "value2", 10)
	if err != nil {
		t.Error(err)
		return
	}
	err = storage.DeleteAll()
	if err != nil {
		t.Error(err)
		return
	}
	value, err := storage.Get("key")
	if err != nil {
		t.Error(err)
		return
	}
	if value != "" {
		t.Error("value error")
		return
	}
	value, err = storage.Get("key2")
	if err != nil {
		t.Error(err)
		return
	}
	if value != "" {
		t.Error("value error")
		return
	}
}

func TestRedisIncrement(t *testing.T) {
	storage, err := NewRedisStorage("127.0.0.1:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}
	err = storage.SetInt("key", 10, 10)
	if err != nil {
		t.Error(err)
		return
	}
	newValue, err := storage.Increment("key", 12)
	if err != nil {
		t.Error(err)
		return
	}

	if newValue != 22 {
		t.Error("value error")
		return
	}
}

func TestRedisDecrement(t *testing.T) {
	storage, err := NewRedisStorage("127.0.0.1:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}
	err = storage.SetInt("key", 10, 10)
	if err != nil {
		t.Error(err)
		return
	}
	newValue, err := storage.Decrement("key", 12)
	if err != nil {
		t.Error(err)
		return
	}

	if newValue != -2 {
		t.Error("value error", newValue)
		return
	}
}
*/
