package safe_stop

import "sync"

type SafeStop struct {
	wg *sync.WaitGroup
}

func NewSafeStop() *SafeStop {
	o := &SafeStop{wg: &sync.WaitGroup{}}
	return o
}

func (s *SafeStop) Add(in int) { s.wg.Add(in) }
func (s *SafeStop) Done()      { s.wg.Done() }
func (s *SafeStop) Wait()      { s.wg.Wait() }

var Default = NewSafeStop()

func Add(in int) { Default.Add(in) }
func Done()      { Default.Done() }
func Wait()      { Default.Wait() }
