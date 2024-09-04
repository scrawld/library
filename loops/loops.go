package loops

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"
)

type Entry struct {
	Spec time.Duration
	Job  func()
}

type Loops struct {
	entries   []*Entry
	running   bool
	runningMu sync.Mutex
	jobWaiter sync.WaitGroup
}

func New() *Loops {
	o := &Loops{}
	return o
}

func (l *Loops) AddFunc(spec time.Duration, cmd func()) {
	l.runningMu.Lock()
	defer l.runningMu.Unlock()

	entry := &Entry{
		Spec: spec,
		Job:  cmd,
	}
	l.entries = append(l.entries, entry)
}

func (l *Loops) Start() {
	l.runningMu.Lock()
	defer l.runningMu.Unlock()

	if l.running {
		return
	}
	l.running = true

	for _, entry := range l.entries {
		l.jobWaiter.Add(1)

		go func(entry *Entry) {
			defer l.jobWaiter.Done()

			for l.running {
				l.runWithRecovery(entry.Job)
				time.Sleep(entry.Spec)
			}
		}(entry)
	}
}

func (l *Loops) runWithRecovery(f func()) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("loops: panic running job: %v\n%s", r, buf)
		}
	}()
	f()
}

func (l *Loops) Stop() context.Context {
	l.runningMu.Lock()
	defer l.runningMu.Unlock()

	if l.running {
		l.running = false
	}

	// Create a new context and cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Start a goroutine to cancel all active consumers
	go func() {
		l.jobWaiter.Wait() // Wait for consumer jobs to complete
		cancel()           // Cancel the context once shutdown is complete
	}()
	return ctx
}
