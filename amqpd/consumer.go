package amqpd

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var consumerSeq uint64

type entry struct {
	Queue   string
	Handler func([]byte) error
}

// AmqpdConsumer is a struct for an AMQP consumer, used for asynchronously consuming messages from multiple queues.
type AmqpdConsumer struct {
	entries   map[string]*entry
	cli       *Amqpd
	running   bool
	runningMu sync.Mutex
	jobWaiter sync.WaitGroup
}

// NewAmqpdConsumer creates a new AmqpdConsumer instance.
func NewAmqpdConsumer() (*AmqpdConsumer, error) {
	cli, err := New()
	if err != nil {
		return nil, fmt.Errorf("amqpd connect err, %s", err)
	}
	return &AmqpdConsumer{
		entries:   make(map[string]*entry),
		cli:       cli,
		running:   false,
		runningMu: sync.Mutex{},
	}, nil
}

// AddFunc adds a queue consumption configuration to the AmqpdConsumer.
func (ac *AmqpdConsumer) AddFunc(queue, consumer string, fn func([]byte) error) {
	ac.runningMu.Lock()
	defer ac.runningMu.Unlock()

	suffix := "-" + strconv.FormatUint(atomic.AddUint64(&consumerSeq, 1), 10)

	ac.entries[consumer+suffix] = &entry{
		Queue:   queue,
		Handler: fn,
	}
	return
}

// Start starts the AmqpdConsumer and begins asynchronous consumption of configured queues.
func (ac *AmqpdConsumer) Start() {
	ac.runningMu.Lock()
	defer ac.runningMu.Unlock()

	if ac.running {
		return
	}
	ac.running = true

	for k, v := range ac.entries {
		ac.jobWaiter.Add(1)

		go func(c string, e *entry) {
			defer ac.jobWaiter.Done()
			ac.run(c, e)
		}(k, v)
	}
	return
}

// run starts an asynchronous consumer for a specified queue.
func (ac *AmqpdConsumer) run(csr string, e *entry) {
	for ac.running {
		err := ac.consume(e.Queue, csr, e.Handler)
		if err != nil {
			log.Printf("amqpd-consumer: run error: %s\n", err)
			time.Sleep(time.Second * 15)
			continue
		}
		if !ac.running {
			break
		}
		time.Sleep(time.Second * 15)
	}
	return
}

// consume connects to the specified queue and handles message consumption.
func (ac *AmqpdConsumer) consume(queue, consumer string, handler func([]byte) error) error {
	deliveries, err := ac.cli.Consume(queue, consumer)
	if err != nil {
		return fmt.Errorf("amqpd consume err, %s", err)
	}
	for dely := range deliveries {
		err := ac.runWithRecovery(handler, dely.Body)
		if err != nil {
			dely.Reject(true)
			continue
		}
		dely.Ack(false)
	}
	return nil
}

// runWithRecovery is a utility method for running a function 'f' with panic recovery.
func (ac *AmqpdConsumer) runWithRecovery(f func([]byte) error, body []byte) error {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("amqpd-consumer: panic running job: %v\n%s\n", r, buf)
		}
	}()
	return f(body)
}

// Stop stops the AmqpdConsumer, which includes canceling all active consumers,
// waiting for the consumer jobs to complete, and closing the AMQP channel.
// It returns a context.Context that is canceled when the AmqpdConsumer has
// completed its shutdown process. Once the channel is closed, this AmqpdConsumer
// cannot be used for further operations.
//
// This method should be called when you want to gracefully shut down the AmqpdConsumer.
//
// Example:
//
//	ctx := amqpdConsumer.Stop()
//	<-ctx.Done() // Wait for the AmqpdConsumer to complete its shutdown.
func (ac *AmqpdConsumer) Stop() context.Context {
	ac.runningMu.Lock()
	defer ac.runningMu.Unlock()

	if ac.running {
		ac.running = false
	}

	// Create a new context and cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Start a goroutine to cancel all active consumers
	go func() {
		for csr := range ac.entries {
			ac.cli.Cancel(csr)
		}
		ac.jobWaiter.Wait() // Wait for consumer jobs to complete
		ac.cli.Close()      // Close the AMQP channel
		cancel()            // Cancel the context once shutdown is complete
	}()
	return ctx
}
