package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Lock struct {
	key string
	tag string
}

func NewLock(key string) *Lock {
	uid, _ := uuid.NewRandom()
	o := &Lock{
		key: fmt.Sprintf("%s.lock.%s", KeyPrefix, key),
		tag: uid.String(),
	}
	return o
}

func (l *Lock) Lock(expire time.Duration) (bool, error) {
	return GetClient().SetNX(context.Background(), l.key, l.tag, expire).Result()
}

func (l *Lock) WaitAndLock(expire time.Duration) (err error) {
	for {
		var ok bool
		ok, err = GetClient().SetNX(context.Background(), l.key, l.tag, expire).Result()
		if err != nil {
			return
		}
		if ok {
			break
		}
		time.Sleep(5 * time.Millisecond) // 5毫秒
	}
	return
}

func (l *Lock) Unlock() (err error) {
	tag, err := GetClient().Get(context.Background(), l.key).Result()
	if err == Nil {
		return nil
	}
	if err != nil {
		err = fmt.Errorf("get %s error, %s", l.key, err)
		return
	}
	if tag != l.tag {
		err = fmt.Errorf("there is no lock")
		return
	}
	return GetClient().Del(context.Background(), l.key).Err()
}
