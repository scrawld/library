package amqpd

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/scrawld/library/config"

	"github.com/streadway/amqp"
)

var (
	mutex      sync.Mutex
	Connection *amqp.Connection
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
		if config.Get().Rabbitmq.TlsProtocols {
			protocol = "amqps"
			amqpConfig.TLSClientConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
		} else {
			protocol = "amqp"
		}
		url := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
			protocol,
			config.Get().Rabbitmq.Username,
			config.Get().Rabbitmq.Password,
			config.Get().Rabbitmq.Host,
			config.Get().Rabbitmq.Port,
			config.Get().Rabbitmq.Vhost,
		)
		Connection, err = amqp.DialConfig(url, amqpConfig)
		if err != nil {
			err = fmt.Errorf("failed amqpd connect %s, %s", err, url)
			return
		}
	}
	if err != nil {
		err = fmt.Errorf("failed amqpd connect %s", err)
		return
	}
	return
}
