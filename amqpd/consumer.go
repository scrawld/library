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

	for csr, entry := range ac.entries {
		go ac.run(csr, entry)
	}
	return
}

// run starts an asynchronous consumer for a specified queue.
func (ac *AmqpdConsumer) run(csr string, e *entry) {
	for ac.running {
		err := ac.consume(e.Queue, csr, e.Handler)
		if err != nil {
			log.Printf("amqpd-consumer: run error: %s\n", err)
			time.Sleep(time.Second * 20)
			continue
		}
		time.Sleep(time.Second * 20)
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
		if !ac.running {
			break
		}
		ac.jobWaiter.Add(1)

		err := ac.runWithRecovery(handler, dely.Body)
		if err != nil {
			dely.Reject(true)
			continue
		}
		dely.Ack(false)

		ac.jobWaiter.Done()
	}
	return nil
}

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
		for csr := range ac.entries {
			ac.cli.Cancel(csr)
		}
		ac.jobWaiter.Wait()
		ac.cli.Close()
		cancel()
	}()
	return ctx
}
