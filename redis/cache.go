package redis

import (
	"context"
	"fmt"
	"time"
)

type Cache struct {
	ctx  context.Context
	pref string
}

func New() (r *Cache) {
	r = &Cache{
		ctx:  context.Background(),
		pref: fmt.Sprintf("%s.store.", KeyPrefix),
	}
	return
}

func (c *Cache) SetPref(pref string) *Cache {
	c.pref = pref
	return c
}

func (c *Cache) EmptyPref() *Cache {
	c.pref = ""
	return c
}

func (c *Cache) Set(key string, val interface{}) error {
	return GetClient().Set(c.ctx, fmt.Sprintf("%s%s", c.pref, key), val, 0).Err()
}

func (c *Cache) SetEX(key string, val interface{}, expire time.Duration) error {
	return GetClient().Set(c.ctx, fmt.Sprintf("%s%s", c.pref, key), val, expire).Err()
}

func (c *Cache) Get(key string) (string, error) {
	return GetClient().Get(c.ctx, fmt.Sprintf("%s%s", c.pref, key)).Result()
}

func (c *Cache) GetInt64(key string) (int64, error) {
	return GetClient().Get(c.ctx, fmt.Sprintf("%s%s", c.pref, key)).Int64()
}

func (c *Cache) Del(keys ...string) error {
	sk := []string{}
	for _, k := range keys {
		sk = append(sk, fmt.Sprintf("%s%s", c.pref, k))
	}
	return GetClient().Del(c.ctx, sk...).Err()
}

func (c *Cache) Exists(key string) (int64, error) {
	return GetClient().Exists(c.ctx, fmt.Sprintf("%s%s", c.pref, key)).Result()
}

func (c *Cache) Expire(key string, expiration time.Duration) (bool, error) {
	return GetClient().Expire(c.ctx, fmt.Sprintf("%s%s", c.pref, key), expiration).Result()
}

func (c *Cache) ExpireAt(key string, tm time.Time) (bool, error) {
	return GetClient().ExpireAt(c.ctx, fmt.Sprintf("%s%s", c.pref, key), tm).Result()
}

func (c *Cache) HGet(key string, field string) (string, error) {
	return GetClient().HGet(c.ctx, fmt.Sprintf("%s%s", c.pref, key), field).Result()
}

func (c *Cache) HSet(key string, field string, val string) error {
	return GetClient().HSet(c.ctx, fmt.Sprintf("%s%s", c.pref, key), field, val).Err()
}

func (c *Cache) HDel(key string, field string) error {
	return GetClient().HDel(c.ctx, fmt.Sprintf("%s%s", c.pref, key), field).Err()
}
