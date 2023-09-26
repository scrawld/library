package amqpd

import "errors"

// Publish publishes a message to the specified exchange with the given routing key
// using the default Amqpd instance (Default).
func Publish(exchange, key string, body []byte) error {
	if Default == nil {
		return errors.New("default amqpd instance is not initialized")
	}
	return Default.Publish(exchange, key, body)
}
