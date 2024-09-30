package amqpd

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Amqpd struct {
	channel *amqp.Channel
	stop    chan struct{}
}

// New creates a new Amqpd instance and initializes its channel.
func New() (*Amqpd, error) {
	ad := &Amqpd{
		stop: make(chan struct{}),
	}
	if err := ad.initChannel(); err != nil {
		return nil, err
	}
	go ad.redial()
	return ad, nil
}

// initChannel initializes the AMQP channel for the Amqpd instance.
func (ad *Amqpd) initChannel() error {
	if Connection == nil || Connection.IsClosed() {
		// amqpd connection
		if err := Init(); err != nil {
			return fmt.Errorf("amqpd connection error: %s", err)
		}
	}
	if ad.channel != nil && !ad.channel.IsClosed() {
		ad.channel.Close()
	}
	// In a situation where Close is not called, there can be up to 2047 simultaneous channels.
	channel, err := Connection.Channel()
	if err != nil {
		return fmt.Errorf("open channel error: %s", err)
	}
	ad.channel = channel
	return nil
}

// redial monitors the channel and re-establishes it if it's closed.
func (ad *Amqpd) redial() {
	printf := func(format string, v ...any) { log.Printf("amqpd-redial: "+format, v...) }
	for {
		select {
		case <-ad.stop:
			return
		case closeErr := <-ad.channel.NotifyClose(make(chan *amqp.Error)):
			printf("channel closing: %s", closeErr)
			for {
				select {
				case <-ad.stop:
					return
				default:
				}
				printf("reconnecting...")
				if err := ad.initChannel(); err != nil {
					printf("reconnect error: %s", err)
					time.Sleep(time.Second * 10)
					continue
				}
				printf("channel re-established")
				break
			}
		}
	}
}

// Cancel stops deliveries to the consumer chan established in Channel.Consume and identified by consumer.
func (ad *Amqpd) Cancel(consumer string) error {
	return ad.channel.Cancel(consumer, false)
}

// Close closes the Amqpd instance's channel and stops the redialing process.
func (ad *Amqpd) Close() error {
	ad.stop <- struct{}{}
	return ad.channel.Close()
}

// DeclareExchange declares an exchange on the AMQP server with the given name and type.
func (ad *Amqpd) ExchangeDeclare(name string, kind string) error {
	return ad.channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

// Publish publishes a message to the specified exchange with the given routing key.
func (ad *Amqpd) Publish(exchange, key string, body []byte) error {
	return ad.channel.Publish(exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body})
}

// QueueDeclare declares a queue with the given name on the AMQP server.
func (ad *Amqpd) QueueDeclare(name string) (amqp.Queue, error) {
	return ad.channel.QueueDeclare(name, true, false, false, false, nil)
}

// QueueBind binds a queue to an exchange with a routing key.
func (ad *Amqpd) QueueBind(name, key, exchange string) error {
	return ad.channel.QueueBind(name, key, exchange, false, nil)
}

// Consume starts consuming messages from a queue identified by its name.
func (ad *Amqpd) Consume(queue, consumer string) (<-chan amqp.Delivery, error) {
	return ad.channel.Consume(queue, consumer, false, false, false, false, nil)
}
