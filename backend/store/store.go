package store

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Store interface {
	Set(key string, value interface{}) error
	Get(key string) ([]interface{}, error)
	Exist(key string) (bool, error)
	Keys(key string) ([]string, error)
	Del(key string)
}

type Redis struct {
	pool *redis.Pool
}

func NewRedisStore(dialFunc func() (redis.Conn, error)) (Store, error) {
	if dialFunc == nil {
		return nil, fmt.Errorf("Nil dial function provided")
	}
	redisPool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   3,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial:        dialFunc,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &Redis{redisPool}, nil
}

// Set wraps Redis HMSET
func (r *Redis) Set(key string, value interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(value)...)
	if err != nil {
		log.Printf("Fail to set for key %s, value: %v", key, value)
	}
	return err
}

// Keys wraps Redis KEYS
func (r *Redis) Keys(key string) ([]string, error) {
	conn := r.pool.Get()
	defer conn.Close()
	value, err := redis.Strings(conn.Do("KEYS", key))
	if err != nil {
		log.Printf("Fail to list key %s, error: %v", key, err)
	}
	return value, err
}

// Get wraps Redis HGETALL
func (r *Redis) Get(key string) ([]interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()
	value, err := redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		log.Printf("Fail to get data with key %s, error: %v", key, err)
	}
	return value, err
}

// Exist wraps Redis EXSTS
func (r *Redis) Exist(key string) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("EXISTS", key))
}

// Del wraps Redis DEL
func (r *Redis) Del(key string) {
	conn := r.pool.Get()
	defer conn.Close()

	keys, _ := r.Keys(key)
	for _, val := range keys {
		conn.Do("DEL", val)
	}
}
