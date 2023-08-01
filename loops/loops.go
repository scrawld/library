package loops

import (
	"log"
	"os"
	"runtime"
	"time"
)

type Entry struct {
	Spec time.Duration
	Job  func()
}

type Loops struct {
	ExeDir  string
	entries []*Entry
	running bool
}

func New() *Loops {
	o := &Loops{
		ExeDir:  os.Getenv("EXECDIR"),
		running: true,
	}
	return o
}

func (this *Loops) AddFunc(spec time.Duration, cmd func()) {
	entry := &Entry{
		Spec: spec,
		Job:  cmd,
	}
	this.entries = append(this.entries, entry)
}

func (this *Loops) Start() {
	for _, entry := range this.entries {
		go func(entry *Entry) {
			for this.running {
				this.runWithRecovery(entry.Job)
				time.Sleep(entry.Spec)
			}
		}(entry)
	}
}

func (this *Loops) runWithRecovery(f func()) {
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

func (this *Loops) Stop() {
	this.running = false
}
