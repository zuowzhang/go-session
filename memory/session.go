package memory

import (
	"time"
	"session"
	"sync"
	"container/list"
	"fmt"
)

var memoryProvider = &MemoryProvider{
	sessions:make(map[string]*list.Element, 0),
	sessionList:list.New(),
}

func init() {
	fmt.Printf("memory provider init")
	session.RegisterProvider(session.PROVIDER_MEMORY, memoryProvider)
}

type MemoryStore struct {
	id               string
	lastTimeAccessed time.Time
	values           map[interface{}]interface{}
}

func (store *MemoryStore)Set(key, value interface{}) error {
	store.values[key] = value
	memoryProvider.update(store.id)
	return nil
}

func (store *MemoryStore)Get(key interface{}) interface{} {
	value, ok := store.values[key]
	if ok {
		return value
	}
	return nil
}

func (store *MemoryStore)Delete(key interface{}) error {
	delete(store.values, key)
	memoryProvider.update(store.id)
	return nil
}

func (store *MemoryStore)Id() string {
	return store.id
}

type MemoryProvider struct {
	lock        sync.RWMutex
	sessions    map[string]*list.Element
	sessionList *list.List
}

func (provider *MemoryProvider)Read(sid string) (session.Session, error) {
	ok := false
	provider.lock.RLock()
	element, ok := provider.sessions[sid]
	provider.lock.RUnlock()
	if !ok {
		provider.lock.Lock()
		defer provider.lock.Unlock()
		element = provider.sessionList.PushFront(&MemoryStore{
			id:sid,
			lastTimeAccessed:time.Now(),
			values:make(map[interface{}]interface{}),
		})
		provider.sessions[sid] = element

	}
	return element.Value.(*MemoryStore), nil
}

func (provider *MemoryProvider)Remove(sid string) error {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	element, ok := provider.sessions[sid]
	if ok {
		delete(provider.sessions, sid)
		provider.sessionList.Remove(element)
	}
	return nil
}

func (provider *MemoryProvider)Gc(maxLifeTime int64) {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	for {
		element := provider.sessionList.Back()
		if element == nil {
			break
		}
		store := element.Value.(*MemoryStore)
		if store.lastTimeAccessed.Unix() + maxLifeTime < time.Now().Unix() {
			provider.sessionList.Remove(element)
			delete(provider.sessions, store.id)
		} else {
			break
		}
	}
}

func (provider *MemoryProvider)update(sid string) {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	element, ok := provider.sessions[sid]
	if ok {
		element.Value.(*MemoryStore).lastTimeAccessed = time.Now()
		provider.sessionList.MoveToFront(element)
	}
}
