//go:build !solution

package keylock

import (
	"sync"
)

type KeyLock struct {
	mu    sync.Mutex
	locks map[string]chan struct{}
}

func New() *KeyLock {
	return &KeyLock{
		locks: make(map[string]chan struct{}),
	}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	l.mu.Lock()

	waitChans := make([]chan struct{}, 0, len(keys))
	for _, key := range keys {
		if ch, ok := l.locks[key]; !ok {
			// No one is currently locking this key, create a channel and assign it
			ch = make(chan struct{})
			l.locks[key] = ch
			close(ch)
		} else {
			waitChans = append(waitChans, ch)
		}
	}

	l.mu.Unlock()

	// Now wait for all the channels to be closed (unlocked)
	for _, ch := range waitChans {
		select {
		case <-ch:
			// continue to wait for the next channel
		case <-cancel:
			// Cancellation requested, return true and no unlock function
			return true, nil
		}
	}

	// All keys are free, proceed to lock them and prepare the unlock function
	unlockChans := make(map[string]chan struct{})
	for _, key := range keys {
		ch := make(chan struct{})
		unlockChans[key] = ch
		l.mu.Lock()
		l.locks[key] = ch
		l.mu.Unlock()
	}

	unlock = func() {
		l.mu.Lock()
		for key, ch := range unlockChans {
			close(ch)
			delete(l.locks, key)
		}
		l.mu.Unlock()
	}

	return false, unlock
}
