package keylock

import "sync"

type KeyLock struct {
	giantLock sync.Mutex
	locks     map[string]*sync.Mutex
}

func NewKeyLock() *KeyLock {
	return &KeyLock{
		giantLock: sync.Mutex{},
		locks:     map[string]*sync.Mutex{},
	}
}

func (self *KeyLock) getLock(key string) *sync.Mutex {
	if lock, ok := self.locks[key]; ok {
		return lock
	}

	self.giantLock.Lock()

	if lock, ok := self.locks[key]; ok {
		self.giantLock.Unlock()
		return lock
	}

	lock := &sync.Mutex{}
	self.locks[key] = lock
	self.giantLock.Unlock()
	return lock
}

func (self *KeyLock) Lock(key string) {
	self.getLock(key).Lock()
}

func (self *KeyLock) Unlock(key string) {
	self.locks[key].Unlock()
}

type RWKeyLock struct {
	giantLock sync.Mutex
	locks     map[string]*sync.RWMutex
}

func NewRWKeyLock() *RWKeyLock {
	return &RWKeyLock{
		giantLock: sync.Mutex{},
		locks:     map[string]*sync.RWMutex{},
	}
}

func (self *RWKeyLock) getLock(key string) *sync.RWMutex {
	if lock, ok := self.locks[key]; ok {
		return lock
	}

	self.giantLock.Lock()

	if lock, ok := self.locks[key]; ok {
		self.giantLock.Unlock()
		return lock
	}

	lock := &sync.RWMutex{}
	self.locks[key] = lock
	self.giantLock.Unlock()
	return lock
}

func (self *RWKeyLock) Lock(key string) {
	self.getLock(key).Lock()
}

func (self *RWKeyLock) Unlock(key string) {
	self.locks[key].Unlock()
}

func (self *RWKeyLock) RLock(key string) {
	self.getLock(key).RLock()
}

func (self *RWKeyLock) RUnlock(key string) {

	self.locks[key].RUnlock()
}
