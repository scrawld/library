package amqpd

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/scrawld/library/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	mutex      sync.Mutex
	Connection *amqp.Connection
	Default    *Amqpd
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
