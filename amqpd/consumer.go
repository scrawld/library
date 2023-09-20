package amqpd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type entry struct {
	Queue    string
	Consumer string
	Handler  func([]byte) error
}

// AmqpdConsumer is a struct for an AMQP consumer, used for asynchronously consuming messages from multiple queues.
type AmqpdConsumer struct {
	entries   []*entry
	clis      []*Amqpd
	running   bool
	runningMu sync.Mutex
	jobWaiter sync.WaitGroup
}

// NewAmqpdConsumer creates a new AmqpdConsumer instance.
func NewAmqpdConsumer() *AmqpdConsumer {
	o := &AmqpdConsumer{
		running:   false,
		runningMu: sync.Mutex{},
	}
	return o
}

// AddFunc adds a queue consumption configuration to the AmqpdConsumer.
func (ac *AmqpdConsumer) AddFunc(queue, consumer string, fn func([]byte) error) {
	ac.runningMu.Lock()
	defer ac.runningMu.Unlock()

	ac.entries = append(ac.entries, &entry{
		Queue:    queue,
		Consumer: consumer,
		Handler:  fn,
	})
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

	for _, entry := range ac.entries {
		go ac.run(entry)
	}
	return
}

// run starts an asynchronous consumer for a specified queue.
func (ac *AmqpdConsumer) run(e *entry) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Fprintf(os.Stderr, "amqpd: panic running job: %v\n%s\n", r, buf)
		}
	}()
	for ac.running {
		err := ac.consume(e.Queue, e.Consumer, e.Handler)
		if err != nil {
			fmt.Fprintf(os.Stderr, "run error: %s\n", err)
			time.Sleep(time.Second * 5)
			return
		}
		time.Sleep(time.Second)
	}
	return
}

// consume connects to the specified queue and handles message consumption.
func (ac *AmqpdConsumer) consume(queue, consumer string, handler func([]byte) error) error {
	cli, err := New()
	if err != nil {
		return fmt.Errorf("amqpd connect err, %s", err)
	}
	ac.clis = append(ac.clis, cli)

	var deliveries <-chan amqp.Delivery
	deliveries, err = cli.Consume(queue, consumer)
	if err != nil {
		return fmt.Errorf("amqpd consume err, %s", err)
	}
	for dely := range deliveries {
		if !ac.running {
			break
		}
		ac.jobWaiter.Add(1)

		err := handler(dely.Body)
		if err != nil {
			dely.Reject(true)
			continue
		}
		dely.Ack(false)

		ac.jobWaiter.Done()
	}
	return nil
}

// Stop stops the AmqpdConsumer, waits for all jobs to complete, and closes AMQP connections.
// A context is returned so the caller can wait for running jobs to complete.
func (ac *AmqpdConsumer) Stop() context.Context {
	ac.runningMu.Lock()
	defer ac.runningMu.Unlock()

	if ac.running {
		ac.running = false
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ac.jobWaiter.Wait()
		for _, cli := range ac.clis {
			cli.Close()
		}
		cancel()
	}()
	return ctx
}
