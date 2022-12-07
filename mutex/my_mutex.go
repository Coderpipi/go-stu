package mutex

import "time"

type Mutex struct {
	ch chan struct{}
}

func NewMutex() *Mutex {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return &Mutex{ch: ch}

}

func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

func (m *Mutex) TryLockWithTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-timer.C:
	}
	return false
}

func (m *Mutex) Lock() {
	<-m.ch
}

func (m *Mutex) UnLock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock for unlocked mutex")
	}
}
