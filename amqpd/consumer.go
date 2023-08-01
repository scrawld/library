package amqpd

import (
	"fmt"
	"os"
	"time"
)

type entry struct {
	Queue    string
	Consumer string
	Handler  func([]byte) error
}

type AmqpdConsumer struct {
	running bool
	entries []*entry
	clis    []*Amqpd
}

func NewAmqpdConsumer() *AmqpdConsumer {
	o := &AmqpdConsumer{
		running: true,
	}
	return o
}

func (this *AmqpdConsumer) AddFunc(queue, consumer string, fn func([]byte) error) {
	this.entries = append(this.entries, &entry{
		Queue:    queue,
		Consumer: consumer,
		Handler:  fn,
	})
	return
}

func (this *AmqpdConsumer) Start() {
	for _, entry := range this.entries {
		go this.run(entry)
	}
	return
}

func (this *AmqpdConsumer) run(e *entry) {
	fn := func(ent *entry) {
		cli, err := New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "amqpd connect err, %s\n", err)
			time.Sleep(time.Second * 5)
			return
		}
		defer cli.Close()

		this.clis = append(this.clis, cli)
		if err = cli.Consume(ent.Queue, ent.Consumer, ent.Handler); err != nil {
			fmt.Fprintf(os.Stderr, "amqpd consume err, %s\n", err)
			time.Sleep(time.Second * 5)
			return
		}
		time.Sleep(time.Second)
	}
	for this.running {
		fn(e)
	}
	return
}

func (this *AmqpdConsumer) Stop() (err error) {
	this.running = false
	for _, cli := range this.clis {
		cli.Close()
	}
	return
}
