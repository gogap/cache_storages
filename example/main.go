package main

import (
	"fmt"

	"github.com/gogap/cache_storages"
)

func main() {
	storage, err := cache_storages.NewMemcachedStorage("127.0.0.1:11211")
	if err != nil {
		panic(err)
	}

	obj := map[string]interface{}{}
	err = storage.GetObject("62e76881-fd48-49fa-65fc-fa4ae43fa8e6", &obj)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(obj)
	}
}
