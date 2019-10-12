package store

import (
	"sync"
)

type Store struct {
	keys *sync.Map
}

// Returns a new, empty store
func StoreService() Store {
	s := Store{}
	s.keys = new(sync.Map)
	return s
}

// Set a key
func (s Store) Set(key string, value interface{}) {
	s.keys.Store(key, value)
}

// Get a key
func (s Store) Get(key string) (interface{}, bool) {
	return s.keys.Load(key)
}

// Delete a key
func (s Store) Delete(key string) {
	s.keys.Delete(key)
}

// Delete all keys
func (s Store) Flush() {
	s.keys = new(sync.Map)
}
