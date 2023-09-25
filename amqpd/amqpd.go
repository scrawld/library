package amqpd

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Amqpd struct {
	channel *amqp.Channel
}

func New() (r *Amqpd, err error) {
	if err = reConnect(); err != nil {
		err = fmt.Errorf("re connect: %s", err)
		return
	}
	var chann *amqp.Channel
	// In a situation where Close is not called, there can be up to 2047 simultaneous channels.
	chann, err = Connection.Channel()
	if err != nil {
		err = fmt.Errorf("failed to open a channel %s", err)
		return
	}
	r = &Amqpd{channel: chann}
	return
}

func reConnect() (err error) {
	for i := 0; i < 5; i++ {
		if Connection != nil && !Connection.IsClosed() {
			return
		}
		err = Init()
		if err == nil {
			return
		}
		err = fmt.Errorf("amqpd connection error: %s", err)
	}
	return
}

// Cancel stops deliveries to the consumer chan established in Channel.Consume and identified by consumer.
func (ad *Amqpd) Cancel(consumer string) error {
	return ad.channel.Cancel(consumer, false)
}

func (ad *Amqpd) Close() error {
	return ad.channel.Close()
}

// ExchangeDeclare
func (ad *Amqpd) ExchangeDeclare(name string, kind ExchangeType) (err error) {
	return ad.channel.ExchangeDeclare(name, string(kind), true, false, false, false, nil)
}

// Publish
func (ad *Amqpd) Publish(exchange, key string, body []byte) (err error) {
	return ad.channel.Publish(exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body})
}

// QueueDeclare
func (ad *Amqpd) QueueDeclare(name string) (err error) {
	_, err = ad.channel.QueueDeclare(name, true, false, false, false, nil)
	return
}

// QueueBind
func (ad *Amqpd) QueueBind(name, key, exchange string) (err error) {
	return ad.channel.QueueBind(name, key, exchange, false, nil)
}

// Consume
func (ad *Amqpd) Consume(queue, consumer string) (<-chan amqp.Delivery, error) {
	return ad.channel.Consume(queue, consumer, false, false, false, false, nil)
}
