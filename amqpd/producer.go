package amqpd

// Publish publishes a message to the specified exchange with the given routing key
// using the default Amqpd instance (Default).
func Publish(exchange, key string, body []byte) error {
	return Default.Publish(exchange, key, body)
}
