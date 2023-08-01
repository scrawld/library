package redis

import (
	"crypto/tls"

	"github.com/scrawld/library/config"

	"github.com/go-redis/redis/v8"
)

var (
	Client    *redis.Client
	KeyPrefix = "keyPrefix" // your project name
	Nil       = redis.Nil
)

func Init() {
	var (
		redisCfg = config.Get().Redis
		opt      = &redis.Options{
			Addr:     redisCfg.Addr,
			Password: redisCfg.Password, // no password set
			DB:       redisCfg.DB,       // use default DB
		}
	)
	if len(redisCfg.Username) != 0 {
		opt.Username = redisCfg.Username
	}
	if redisCfg.TlsProtocols {
		opt.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	Client = redis.NewClient(opt)
}

func GetClient() *redis.Client {
	if Client == nil {
		Init()
	}
	return Client
}
