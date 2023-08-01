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

func (this *Amqpd) Close() error {
	return this.channel.Close()
}

// ExchangeDeclare
func (this *Amqpd) ExchangeDeclare(name string, kind ExchangeType) (err error) {
	return this.channel.ExchangeDeclare(name, string(kind), true, false, false, false, nil)
}

// Publish
func (this *Amqpd) Publish(exchange, key string, body []byte) (err error) {
	return this.channel.Publish(exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body})
}

// QueueDeclare
func (this *Amqpd) QueueDeclare(name string) (err error) {
	_, err = this.channel.QueueDeclare(name, true, false, false, false, nil)
	return
}

// QueueBind
func (this *Amqpd) QueueBind(name, key, exchange string) (err error) {
	return this.channel.QueueBind(name, key, exchange, false, nil)
}

// Consume
func (this *Amqpd) Consume(queue, consumer string, handler func([]byte) error) (err error) {
	var msgs <-chan amqp.Delivery
	msgs, err = this.channel.Consume(queue, consumer, false, false, false, false, nil)
	if err != nil {
		return
	}
	for dely := range msgs {
		err := handler(dely.Body)
		if err != nil {
			dely.Reject(true)
			continue
		}
		dely.Ack(false)
	}
	return
}
