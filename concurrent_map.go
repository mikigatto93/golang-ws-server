// concurrent_map
package main

import "sync"

type ConcurrentMap struct {
	sync.RWMutex
	items map[string]*WSClient
}

func NewConcurrentMap() *ConcurrentMap {
	cm := new(ConcurrentMap)
	cm.items = make(map[string]*WSClient)
	return cm
}

func (cm *ConcurrentMap) Add(key string, item *WSClient) {
	cm.Lock()
	defer cm.Unlock()
	cm.items[key] = item

}

func (cm *ConcurrentMap) Get(key string) (*WSClient, bool) {
	cm.Lock()
	defer cm.Unlock()

	val, ok := cm.items[key]

	return val, ok
}

func (cm *ConcurrentMap) Delete(key string) {
	cm.Lock()
	defer cm.Unlock()

	delete(cm.items, key)
}

func (cm *ConcurrentMap) Iterate(Iterfunc func(key string, value *WSClient)) {
	cm.Lock()
	defer cm.Unlock()

	for key, value := range cm.items {
		Iterfunc(key, value)
	}
}
