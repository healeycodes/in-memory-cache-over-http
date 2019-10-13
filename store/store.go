package store

import (
	"strconv"
	"sync"
)

// Store is an in-memory key/value string database.
type Store struct {
	mutex *sync.Mutex
	store map[string]string
}

// Service returns an empty store limited by a number of keys (use zero for no limit).
func Service(size int) *Store {
	s := &Store{}
	s.mutex = &sync.Mutex{}
	s.store = make(map[string]string)
	return s
}

// Set a key.
func (s *Store) Set(key string, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store[key] = value
}

// Get a key.
func (s *Store) Get(key string) (string, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	current, exist := s.store[key]
	if exist {
		return current, true
	}
	return "", false
}

// Delete a key.
func (s *Store) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.store, key)
}

// CheckAndSet a key. Sets only if the compare matches. Set the key if it doesn't exist.
func (s *Store) CheckAndSet(key string, value string, compare string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	current, exist := s.store[key]
	if !exist || current == compare {
		s.store[key] = value
		return true
	}
	return false
}

// Increment a key by an amount. If doesn't exist, set to amount.
func (s *Store) Increment(key string, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	current, exist := s.store[key]
	if !exist {
		s.store[key] = value
		return nil
	}

	y, err := strconv.Atoi(value)
	x, err := strconv.Atoi(current)
	if err == nil {
		s.store[key] = strconv.Itoa(x + y)
		return nil
	}
	return err
}

// Decrement a key by an amount. If doesn't exist, set to amount.
func (s *Store) Decrement(key string, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	current, exist := s.store[key]
	if !exist {
		s.store[key] = value
		return nil
	}

	y, err := strconv.Atoi(value)
	x, err := strconv.Atoi(current)
	if err == nil {
		s.store[key] = strconv.Itoa(x - y)
		return nil
	}
	return err
}

// Append to a key.
func (s *Store) Append(key string, value string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	current, exists := s.store[key]
	if !exists {
		return false
	}
	s.store[key] = current + value
	return true
}

// Prepend to a key.
func (s *Store) Prepend(key string, value string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	current, exists := s.store[key]
	if !exists {
		return false
	}
	s.store[key] = value + current
	return true
}

// Flush all keys
func (s *Store) Flush() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store = make(map[string]string)
}
