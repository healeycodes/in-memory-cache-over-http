package cache

import (
	"container/list"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// Store contains an LRU Cache
type Store struct {
	Mutex *sync.Mutex
	store map[string]*list.Element
	ll    *list.List
	max   int // Zero for unlimited
}

// Node maps a value to a key
type Node struct {
	key    string
	value  string
	expire int // Unix time
}

// Service returns an empty store
func Service(max int) *Store {
	s := &Store{
		Mutex: &sync.Mutex{},
		store: make(map[string]*list.Element),
		ll:    list.New(),
		max:   max,
	}
	return s
}

// Set a key
func (s *Store) Set(key string, value string, expire int) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.set(key, value, expire)
}

// Internal set
func (s *Store) set(key string, value string, expire int) {
	current, exist := s.store[key]
	if exist != true {
		s.store[key] = s.ll.PushFront(&Node{
			key:    key,
			value:  value,
			expire: expire,
		})
		if s.max != 0 && s.ll.Len() > s.max {
			s.delete(s.ll.Remove(s.ll.Back()).(*Node).key)
		}
		return
	}
	current.Value.(*Node).value = value
	current.Value.(*Node).expire = expire
	s.ll.MoveToFront(current)
}

// Get a key
func (s *Store) Get(key string) (string, bool) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	current, exist := s.store[key]
	if exist {
		expire := int64(current.Value.(*Node).expire)
		if expire == 0 || expire > time.Now().Unix() {
			s.ll.MoveToFront(current)
			return current.Value.(*Node).value, true
		}
	}
	return "", false
}

// Delete an item
func (s *Store) Delete(key string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.delete(key)
}

// Internal delete
func (s *Store) delete(key string) {
	current, exist := s.store[key]
	if exist != true {
		return
	}
	s.ll.Remove(current)
	delete(s.store, key)
}

// CheckAndSet a key. Sets only if the compare matches. Set the key if it doesn't exist
func (s *Store) CheckAndSet(key string, value string, expire int, compare string) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	current, exist := s.store[key]
	if !exist || current.Value.(*Node).value == compare {
		s.set(key, value, expire)
		return true
	}
	return false
}

// Increment a key by an amount. Both value and amount should be integers. If doesn't exist, set to amount
func (s *Store) Increment(key string, value string, expire int) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	current, exist := s.store[key]
	if !exist {
		s.set(key, value, expire)
	}

	y, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	x, err := strconv.Atoi(current.Value.(*Node).value)
	if err != nil {
		return err
	}
	s.set(key, strconv.Itoa(x+y), expire)
	return nil
}

// Decrement a key by an amount. Both value and amount should be integers. If doesn't exist, set to amount
func (s *Store) Decrement(key string, value string, expire int) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	current, exist := s.store[key]
	if !exist {
		s.set(key, value, expire)
	}

	y, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	x, err := strconv.Atoi(current.Value.(*Node).value)
	if err != nil {
		return err
	}
	s.set(key, strconv.Itoa(x-y), expire)
	return nil
}

// Append to a key
func (s *Store) Append(key string, value string, expire int) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	current, exist := s.store[key]
	if !exist {
		s.set(key, value, expire)
		return
	}
	s.set(key, current.Value.(*Node).value+value, expire)
}

// Prepend to a key
func (s *Store) Prepend(key string, value string, expire int) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	current, exist := s.store[key]
	if !exist {
		s.set(key, value, expire)
		return
	}
	s.set(key, value+current.Value.(*Node).value, expire)
}

// Flush all keys
func (s *Store) Flush() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.store = make(map[string]*list.Element)
	s.ll = list.New()
}

// Stats returns up-to-date information about the cache
func (s *Store) Stats() string {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	// TODO (healeycodes)
	// Use json package here
	return fmt.Sprintf(`{"keyCount": %v, "maxSize": %v}`, len(s.store), s.max)
}
