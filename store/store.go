package store

import (
	"strconv"
	"sync"
)

// Store is an in-memory key/value string database.
type Store struct {
	store *sync.Map
}

// Service returns an empty
func Service() Store {
	s := Store{}
	s.store = new(sync.Map)
	return s
}

// Set a key
func (s Store) Set(key string, value string) {
	s.store.Store(key, value)
}

// Get a key
func (s Store) Get(key string) (string, bool) {
	val, exist := s.store.Load(key)
	if exist {
		return string(val.(string)), true
	}
	return "", false
}

// Delete a key
func (s Store) Delete(key string) {
	s.store.Delete(key)
}

// Increment a key. If doesn't exist, increment from zero.
func (s Store) Increment(key string) error {
	val, exist := s.Get(key)
	if !exist {
		s.Set(key, "1")
		return nil
	}

	n, err := strconv.Atoi(val)
	if err == nil {
		s.Set(key, strconv.Itoa(n+1))
		return nil
	}
	return err
}

// Decrement a key. If doesn't exist, decrement from zero.
func (s Store) Decrement(key string) error {
	val, exist := s.Get(key)
	if !exist {
		s.Set(key, "-1")
		return nil
	}

	n, err := strconv.Atoi(val)
	if err == nil {
		s.Set(key, strconv.Itoa(n-1))
		return nil
	}
	return err
}

// Append to a key.
func (s Store) Append(key string, value string) bool {
	current, exists := s.Get(key)
	if !exists {
		return false
	}
	s.Set(key, current+value)
	return true
}
