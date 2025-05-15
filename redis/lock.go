package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Lock struct {
	prefix string
	rawKey string
	tag    string
}

// NewLock 创建一个新的分布式锁实例，使用默认的前缀格式：<KeyPrefix>.lock.<rawKey>
func NewLock(key string) *Lock {
	uid, _ := uuid.NewRandom()
	o := &Lock{
		prefix: fmt.Sprintf("%s.lock.", KeyPrefix),
		rawKey: key,
		tag:    uid.String(),
	}
	return o
}

// SetPrefix 设置锁的前缀
func (l *Lock) SetPrefix(prefix string) *Lock {
	l.prefix = prefix
	return l
}

// EmptyPrefix 清空前缀
func (l *Lock) EmptyPrefix() *Lock {
	l.prefix = ""
	return l
}

// FullKey 返回拼接后的完整 Redis key：<prefix>.<rawKey>
func (l *Lock) FullKey() string {
	return l.prefix + l.rawKey
}

// Lock 尝试获取锁，设置过期时间。如果锁已存在，返回 false
func (l *Lock) Lock(expire time.Duration) (bool, error) {
	return GetClient().SetNX(context.Background(), l.FullKey(), l.tag, expire).Result()
}

// WaitAndLock 阻塞等待直到成功获取锁
func (l *Lock) WaitAndLock(expire time.Duration) (err error) {
	for {
		var ok bool
		ok, err = GetClient().SetNX(context.Background(), l.FullKey(), l.tag, expire).Result()
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

// Unlock 释放锁，仅当当前锁的 tag 与设置时一致才会删除
func (l *Lock) Unlock() (err error) {
	tag, err := GetClient().Get(context.Background(), l.FullKey()).Result()
	if err == Nil {
		return nil
	}
	if err != nil {
		err = fmt.Errorf("get %s error, %s", l.FullKey(), err)
		return
	}
	if tag != l.tag {
		err = fmt.Errorf("there is no lock")
		return
	}
	return GetClient().Del(context.Background(), l.FullKey()).Err()
}
