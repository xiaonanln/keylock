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

func (self *KeyLock) KeyLocker(key string) sync.Locker {
	return self.getLock(key)
}

type KeyRWLock struct {
	giantLock sync.Mutex
	locks     map[string]*sync.RWMutex
}

func NewKeyRWLock() *KeyRWLock {
	return &KeyRWLock{
		giantLock: sync.Mutex{},
		locks:     map[string]*sync.RWMutex{},
	}
}

func (self *KeyRWLock) getLock(key string) *sync.RWMutex {
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

func (self *KeyRWLock) Lock(key string) {
	self.getLock(key).Lock()
}

func (self *KeyRWLock) Unlock(key string) {
	self.locks[key].Unlock()
}

func (self *KeyRWLock) RLock(key string) {
	self.getLock(key).RLock()
}

func (self *KeyRWLock) RUnlock(key string) {
	self.locks[key].RUnlock()
}

func (self *KeyRWLock) KeyLocker(key string) sync.Locker {
	return self.getLock(key)
}

func (self *KeyRWLock) KeyRLocker(key string) sync.Locker {
	return self.getLock(key).RLocker()
}
