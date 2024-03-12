package amqpd

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Host         string `json:"host"`         // RabbitMQ 主机地址
	Port         string `json:"port"`         // RabbitMQ 端口号
	Username     string `json:"username"`     // RabbitMQ 用户名
	Password     string `json:"password"`     // RabbitMQ 密码
	Vhost        string `json:"vhost"`        // RabbitMQ 虚拟主机
	TlsProtocols bool   `json:"tlsProtocols"` // 是否启用 TLS 协议
}

var (
	GlobalConfig Config // 全局配置变量
	mutex        sync.Mutex
	Connection   *amqp.Connection
	Default      *Amqpd
)

// Init initializes the AMQP connection to RabbitMQ.
func Init() (err error) {
	mutex.Lock()
	defer mutex.Unlock()

	if Connection == nil || Connection.IsClosed() {
		var (
			amqpConfig = amqp.Config{
				Heartbeat: 10 * time.Second,
				Locale:    "en_US",
			}
			protocol string
		)
		if GlobalConfig.TlsProtocols {
			protocol = "amqps"
			amqpConfig.TLSClientConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
		} else {
			protocol = "amqp"
		}
		url := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
			protocol,
			GlobalConfig.Username,
			GlobalConfig.Password,
			GlobalConfig.Host,
			GlobalConfig.Port,
			GlobalConfig.Vhost,
		)
		Connection, err = amqp.DialConfig(url, amqpConfig)
		if err != nil {
			err = fmt.Errorf("amqp dial error: %s, %s", err, url)
			return
		}
	}
	if Default == nil {
		Default, err = New()
		if err != nil {
			err = fmt.Errorf("open default channel error: %s", err)
			return
		}
	}
	return
}

// Close closes the AMQP connections and channels if they are non-nil.
func Close() error {
	mutex.Lock()
	defer mutex.Unlock()

	if Default != nil {
		if err := Default.Close(); err != nil {
			return err
		}
	}
	if Connection != nil {
		if err := Connection.Close(); err != nil {
			return err
		}
	}
	return nil
}
