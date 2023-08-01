package amqpd

import (
	"fmt"
	"sync"

	"github.com/scrawld/library/config"

	"github.com/streadway/amqp"
)

var (
	mutex      sync.Mutex
	Connection *amqp.Connection
)

func Init() (err error) {
	mutex.Lock()
	defer mutex.Unlock()

	if Connection == nil || Connection.IsClosed() {
		Connection, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
			config.Get().Rabbitmq.Username,
			config.Get().Rabbitmq.Password,
			config.Get().Rabbitmq.Host,
			config.Get().Rabbitmq.Port,
			config.Get().Rabbitmq.Vhost,
		))
	}
	if err != nil {
		err = fmt.Errorf("failed amqpd connect %s", err)
		return
	}
	return
}
