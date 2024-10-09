package loops

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"
)

// Entry represents a job that will be executed periodically.
// Spec specifies the time interval between job executions.
type Entry struct {
	Spec time.Duration
	Job  func()
}

// Loops is a struct that manages multiple periodic jobs.
type Loops struct {
	ctx    context.Context
	cancel context.CancelFunc

	entries   []*Entry       // List of jobs to be executed
	jobWaiter sync.WaitGroup // WaitGroup to wait for all jobs to complete when stopping
}

// New creates and returns a new Loops instance.
func New() *Loops {
	ctx, cancel := context.WithCancel(context.Background())
	return &Loops{
		ctx:    ctx,
		cancel: cancel,
	}
}

// AddFunc adds a new job to the Loops with the specified interval (Spec) and job function (cmd).
func (l *Loops) AddFunc(spec time.Duration, cmd func()) {
	entry := &Entry{
		Spec: spec,
		Job:  cmd,
	}
	l.entries = append(l.entries, entry)
}

// Start begins executing all added jobs periodically.
// Each job will run in its own goroutine and will continue to run until Stop is called.
func (l *Loops) Start() {
	for _, entry := range l.entries {
		l.jobWaiter.Add(1)

		go l.startJob(entry)
	}
}

// startJob executes a job and handles recovery from panics.
func (l *Loops) startJob(entry *Entry) {
	defer l.jobWaiter.Done()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-time.After(entry.Spec):
			l.runWithRecovery(entry.Job)
		}
	}
}

// runWithRecovery runs the provided job function (f) and recovers from any panic that occurs,
// logging the panic information for debugging purposes.
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

// Stop signals all running jobs to stop and waits for them to complete.
// It returns a context that will be canceled once all jobs have finished.
func (l *Loops) Stop() context.Context {
	// Create a new context and cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Start a goroutine to cancel all active consumers
	go func() {
		l.jobWaiter.Wait() // Wait for consumer jobs to complete
		cancel()           // Cancel the context once shutdown is complete
	}()

	l.cancel() // Signal all jobs to stop
	return ctx
}
