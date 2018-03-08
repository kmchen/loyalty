package store

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

var redisStore, err = NewRedisStore(
	func() (redis.Conn, error) { return redis.Dial("tcp", ":6379") })

func TestKeys(t *testing.T) {
	var key = "testKey1:1234:user"
	testData := &struct {
		A string
		B int
	}{
		A: "a",
		B: 1234,
	}
	var storedKey []string
	var err error
	if err := redisStore.Set(key, testData); err != nil {
		t.Errorf("[Fail] error while setting to Redis %v", err)
	}
	if storedKey, err = redisStore.Keys(key); err != nil {
		t.Errorf("[Fail] error while getting keys from Redis with key %s", err, key)
	}
	if len(storedKey) != 1 && storedKey[0] != key {
		t.Errorf("[Fail] wrong key %s", key)
	}

	redisStore.Del(key)
}

func TestSetGet(t *testing.T) {
	testData := &struct {
		A string
		B int
	}{
		A: "a",
		B: 1234,
	}
	var key = "testKey"
	if err := redisStore.Set(key, testData); err != nil {
		t.Errorf("[Fail] error while setting key to Redis %s, %v", key, err)
	}

	value, err := redisStore.Get(key)
	if err != nil {
		t.Errorf("[Fail] error while getting data from Redis %v", err)
	}
	newData := &struct {
		A string
		B int
	}{}
	if err := redis.ScanStruct(value, newData); err != nil {
		t.Errorf("[Fail] Fail to convert data to struct %v", err)
	}
	if newData.A != "a" && newData.B != 1234 {
		t.Errorf("[Fail] wrong data from Redis")
	}
	redisStore.Del(key)
}
