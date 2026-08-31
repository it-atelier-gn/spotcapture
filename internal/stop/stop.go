package stop

import "sync"

type Stopper struct {
	ch   chan struct{}
	once sync.Once
}

func New() *Stopper {
	return &Stopper{ch: make(chan struct{})}
}

func (s *Stopper) Stop() {
	s.once.Do(func() { close(s.ch) })
}

func (s *Stopper) Done() <-chan struct{} {
	return s.ch
}
