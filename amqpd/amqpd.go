package amqpd

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Amqpd struct {
	channel *amqp.Channel
	stop    chan struct{}
}

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

// initChannel initializes the AMQP channel.
func (ad *Amqpd) initChannel() error {
	if Connection == nil {
		return errors.New("amqpd not initialized")
	}
	if Connection.IsClosed() {
		// amqpd connection
		if err := Init(); err != nil {
			return fmt.Errorf("amqpd connection error: %s", err)
		}
	}
	if ad.channel != nil {
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
			printf("stop")
			return
		case closeErr := <-ad.channel.NotifyClose(make(chan *amqp.Error)):
			printf("channel closing: %s", closeErr)
			for {
				if err := ad.initChannel(); err != nil {
					printf("init channel error: %s", err)
					time.Sleep(time.Second * 10)
					continue
				}
				printf("channel success")
				break
			}
		}
	}
}

// Cancel stops deliveries to the consumer chan established in Channel.Consume and identified by consumer.
func (ad *Amqpd) Cancel(consumer string) error {
	return ad.channel.Cancel(consumer, false)
}

func (ad *Amqpd) Close() error {
	ad.stop <- struct{}{}
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
