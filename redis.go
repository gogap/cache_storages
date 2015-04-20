package cache_storages

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	MAX_IDLE     = 3
	IDLE_TIMEOUT = 180
)

type RedisStorage struct {
	pool  *redis.Pool
	conn  string
	index int
	auth  string
}

func NewAuthRedisStorage(conn string, index int, auth string) (storage *RedisStorage, err error) {
	storage = newRedisStorage(conn, index)
	storage.auth = auth
	storage.init()
	c := storage.pool.Get()
	defer c.Close()
	err = c.Err()
	return
}

func NewRedisStorage(conn string, index int) (storage *RedisStorage, err error) {
	storage = newRedisStorage(conn, index)
	storage.init()
	c := storage.pool.Get()
	defer c.Close()
	err = c.Err()
	return
}

func newRedisStorage(conn string, index int) (storage *RedisStorage) {
	storage = new(RedisStorage)
	storage.conn = conn
	storage.index = index
	return
}

func (p *RedisStorage) StorageType() string {
	return "redis"
}

func (p *RedisStorage) init() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", p.conn)
		if err != nil {
			return
		}
		if p.auth != "" {
			_, err = c.Do("AUTH", p.auth)
		}
		_, selecterr := c.Do("SELECT", p.index)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	p.pool = &redis.Pool{
		MaxIdle:     MAX_IDLE,
		IdleTimeout: IDLE_TIMEOUT * time.Second,
		Dial:        dialFunc,
	}
}

// actually do the redis cmds
func (p *RedisStorage) do(cmd string, args ...interface{}) (interface{}, error) {
	c := p.pool.Get()
	defer c.Close()
	return c.Do(cmd, args...)
}

// put cache to redis
func (p *RedisStorage) SetObject(key string, value interface{}, seconds int32) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	if seconds > 0 {
		_, err = p.do("SETEX", key, seconds, data)
	} else {
		_, err = p.do("SET", key, data)
	}
	return
}

// get cache from redis.
func (p *RedisStorage) GetObject(key string, value interface{}) (err error) {
	vi, err := p.do("GET", key)
	if err != nil || vi == nil {
		return
	}
	data, err := redis.Bytes(vi, err)
	if err != nil {
		return
	}
	return json.Unmarshal(data, value)
}

func (p *RedisStorage) GetMultiObject(keyValues map[string]interface{}) (err error) {
	for key, value := range keyValues {
		if err = p.GetObject(key, value); err != nil {
			return
		}
	}
	return
}

func (p *RedisStorage) Get(key string) (value string, err error) {
	data, err := p.do("GET", key)
	if err != nil || data == nil {
		return
	}
	return redis.String(data, err)
}

func (p *RedisStorage) Set(key, value string, seconds int32) (err error) {
	if seconds > 0 {
		_, err = p.do("SETEX", key, seconds, value)
	} else {
		_, err = p.do("SET", key, value)
	}
	return
}

func (p *RedisStorage) SetInt(key string, value int64, seconds int32) (err error) {
	if seconds > 0 {
		_, err = p.do("SETEX", key, seconds, value)
	} else {
		_, err = p.do("SET", key, value)
	}
	return
}

func (p *RedisStorage) GetInt(key string) (value int64, err error) {
	data, err := p.do("GET", key)
	if err != nil || data == nil {
		return
	}
	return redis.Int64(data, err)
}

func (p *RedisStorage) GetMulti(keys []string) (values map[string]string, err error) {
	values = make(map[string]string)
	for _, key := range keys {
		value, err := p.Get(key)
		if err != nil {
			return values, err
		}
		values[key] = value
	}
	return
}

func (p *RedisStorage) Touch(key string, seconds int32) (err error) {
	_, err = p.do("expire", key, seconds)
	return
}

func (p *RedisStorage) Increment(key string, delta uint64) (newValue int64, err error) {
	return redis.Int64(p.do("INCRBY", key, delta))
}

func (p *RedisStorage) Decrement(key string, delta uint64) (newValue int64, err error) {
	return redis.Int64(p.do("DECRBY", key, delta))
}

func (p *RedisStorage) Delete(key string) (err error) {
	_, err = p.do("DEL", key)
	return
}

func (p *RedisStorage) DeleteAll() (err error) {
	fields, err := redis.Strings(p.do("KEYS", "*"))
	if err != nil {
		return
	}
	for _, field := range fields {
		if _, err = p.do("DEL", field); err != nil {
			return
		}
	}
	return
}
