package amqpd

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// DefaultExchange is the default direct exchange that binds every queue by its
// name. Applications can route to a queue using the queue name as routing key.
const DefaultExchange = amqp.DefaultExchange

// Constants for standard AMQP 0-9-1 exchange types.
const (
	ExchangeDirect  = amqp.ExchangeDirect
	ExchangeFanout  = amqp.ExchangeFanout
	ExchangeTopic   = amqp.ExchangeTopic
	ExchangeHeaders = amqp.ExchangeHeaders
)
