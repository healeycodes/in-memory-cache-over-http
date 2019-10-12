package store

import (
	"strconv"
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
func (s Store) Set(key string, value string) {
	s.keys.Store(key, value)
}

// Get a key
func (s Store) Get(key string) (string, bool) {
	val, exist := s.keys.Load(key)
	return string(val.(string)), exist
}

// Delete a key
func (s Store) Delete(key string) {
	s.keys.Delete(key)
}

// Increment a key. If doesn't exist, increment from zero.
func (s Store) Increment(key string) error {
	val, exist := s.Get(key)
	if !exist {
		s.Set(key, "1")
		return nil
	}
	if n, err := strconv.Atoi(val); err == nil {
		s.Set(key, strconv.Itoa(n+1))
		return nil
	} else {
		return err
	}
}

// Decrement a key. If doesn't exist, decrement from zero.
func (s Store) Decrement(key string) error {
	val, exist := s.Get(key)
	if !exist {
		s.Set(key, "-1")
		return nil
	}
	if n, err := strconv.Atoi(val); err == nil {
		s.Set(key, strconv.Itoa(n-1))
		return nil
	} else {
		return err
	}
}

// Delete all keys
func (s Store) Flush() {
	s.keys = new(sync.Map)
}
