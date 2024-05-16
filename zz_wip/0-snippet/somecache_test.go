package snippet

import (
	"fmt"
	"sync"
	"testing"
)

type SomeCache interface {
	Get(string) (interface{}, bool)
	Set(string, interface{})
}

// simpleCache store something, it implements SomeCache interface
type simpleCache struct {
	lock  sync.RWMutex
	store map[string]interface{}
}

func (s *simpleCache) Get(key string) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	v, ok := s.store[key]
	fmt.Printf("got %q from cache: %t\n", key, ok)
	return v, ok
}

func (s *simpleCache) Set(key string, data interface{}) {
	s.lock.Lock()
	s.store[key] = data
	s.lock.Unlock()
}

var someCache SomeCache

func init() {
	someCache = initSomeCache()
}

func initSomeCache() SomeCache {
	// pointer receiver of simpleCache implements the methods of SomeCache interface
	return &simpleCache{
		store: map[string]interface{}{
			"key1": 0,
			"key2": 1,
		},
	}
}

func Test_SomeCache(t *testing.T) {
	someCache.Set("key3", 100)

	if v, ok := someCache.Get("key2"); !ok {
		fmt.Println("does not exist key2")
	} else {
		fmt.Println("key2 value", v)
	}

	if v, ok := someCache.Get("key4"); !ok {
		fmt.Println("does not exist key4")
	} else {
		fmt.Println("key4 value", v)
	}

}
