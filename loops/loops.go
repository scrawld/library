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
	entries   []*Entry       // List of jobs to be executed
	running   bool           // Flag indicating whether the Loops is running
	runningMu sync.Mutex     // Mutex to protect access to the running flag and entries
	jobWaiter sync.WaitGroup // WaitGroup to wait for all jobs to complete when stopping
}

// New creates and returns a new Loops instance.
func New() *Loops {
	o := &Loops{}
	return o
}

// AddFunc adds a new job to the Loops with the specified interval (Spec) and job function (cmd).
func (l *Loops) AddFunc(spec time.Duration, cmd func()) {
	l.runningMu.Lock()
	defer l.runningMu.Unlock()

	entry := &Entry{
		Spec: spec,
		Job:  cmd,
	}
	l.entries = append(l.entries, entry)
}

// Start begins executing all added jobs periodically.
// Each job will run in its own goroutine and will continue to run until Stop is called.
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
