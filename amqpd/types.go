package amqpd

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type ExchangeType string

// Constants for standard AMQP 0-9-1 exchange types.
const (
	ExchangeDirect  ExchangeType = amqp.ExchangeDirect
	ExchangeFanout  ExchangeType = amqp.ExchangeFanout
	ExchangeTopic   ExchangeType = amqp.ExchangeTopic
	ExchangeHeaders ExchangeType = amqp.ExchangeHeaders
)
