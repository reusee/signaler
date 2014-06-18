package signaler

import (
	"fmt"
	"sync"
)

type Signaler struct {
	cbs   map[string][]interface{}
	calls chan func()
}

func NewSignaler() *Signaler {
	s := &Signaler{
		cbs:   make(map[string][]interface{}),
		calls: make(chan func()),
	}
	go func() {
		for f := range s.calls {
			if f == nil {
				return
			}
			f()
		}
	}()
	return s
}

func (s *Signaler) CloseSignaler() {
	s.calls <- nil
}

func (s *Signaler) OnSignal(signal string, f interface{}) {
	s.calls <- func() {
		s.cbs[signal] = append(s.cbs[signal], f)
	}
}

func (s *Signaler) Signal(signal string, args ...interface{}) {
	s.calls <- func() {
		i := 0
		cbs := s.cbs[signal]
		for i < len(cbs) {
			f := cbs[i]
			switch fun := f.(type) {
			case func():
				fun()
				i++
			case func(interface{}):
				fun(args[0])
				i++
			case func(...interface{}):
				fun(args...)
				i++
			case func() bool:
				if fun() { // one shot
					cbs = append(cbs[:i], cbs[i+1:]...)
				} else {
					i++
				}
			default:
				panic(fmt.Sprintf("invalid signal hadler %T", f))
			}
		}
		s.cbs[signal] = cbs
	}
}

func (s *Signaler) SignalSynced(signal string, args ...interface{}) {
	var lock sync.Mutex
	lock.Lock()
	s.Signal(signal, args...)
	s.calls <- func() {
		lock.Unlock()
	}
	lock.Lock()
}
