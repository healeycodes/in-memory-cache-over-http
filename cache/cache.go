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
	value   string
	expire  int  // Unix time
	deleted bool // Whether this should be cleaned up
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
	current, exist := s.store[key]
	if exist != true {
		s.store[key] = s.ll.PushFront(&Node{
			value:  value,
			expire: expire,
		})
		if s.max != 0 && s.ll.Len() > s.max {
			toBeZeroed := s.ll.Remove(s.ll.Back()).(*Node)
			toBeZeroed.value = ""
			toBeZeroed.deleted = true
		}
		return
	}
	current.Value.(*Node).value = value
	current.Value.(*Node).expire = expire
	s.ll.MoveToFront(current)
}

// Get a key
func (s *Store) Get(key string) (string, bool) {
	current, exist := s.store[key]
	if exist {
		expire := int64(current.Value.(*Node).expire)
		if current.Value.(*Node).deleted != true && (expire == 0 || expire > time.Now().Unix()) {
			s.ll.MoveToFront(current)
			return current.Value.(*Node).value, true
		}
		s.Delete(key) // Clean up item
	}
	return "", false
}

// Delete an item
func (s *Store) Delete(key string) {
	current, exist := s.store[key]
	if exist != true {
		return
	}
	s.ll.Remove(current)
	delete(s.store, key)
}

// CheckAndSet a key. Sets only if the compare matches. Set the key if it doesn't exist
func (s *Store) CheckAndSet(key string, value string, expire int, compare string) bool {
	current, exist := s.store[key]
	if !exist || current.Value.(*Node).value == compare {
		s.Set(key, value, expire)
		return true
	}
	return false
}

// Increment a key by an amount. Both value and amount should be integers. If doesn't exist, set to amount
func (s *Store) Increment(key string, value string, expire int) error {
	current, exist := s.store[key]
	if !exist {
		s.Set(key, value, expire)
	}

	y, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	x, err := strconv.Atoi(current.Value.(*Node).value)
	if err != nil {
		return err
	}
	s.Set(key, strconv.Itoa(x+y), expire)
	return nil
}

// Decrement a key by an amount. Both value and amount should be integers. If doesn't exist, set to amount
func (s *Store) Decrement(key string, value string, expire int) error {
	current, exist := s.store[key]
	if !exist {
		s.Set(key, value, expire)
	}

	y, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	x, err := strconv.Atoi(current.Value.(*Node).value)
	if err != nil {
		return err
	}
	s.Set(key, strconv.Itoa(x-y), expire)
	return nil
}

// Append to a key
func (s *Store) Append(key string, value string, expire int) {
	current, exist := s.store[key]
	if !exist {
		s.Set(key, value, expire)
		return
	}
	s.Set(key, current.Value.(*Node).value+value, expire)
}

// Prepend to a key
func (s *Store) Prepend(key string, value string, expire int) {
	current, exist := s.store[key]
	if !exist {
		s.Set(key, value, expire)
		return
	}
	s.Set(key, value+current.Value.(*Node).value, expire)
}

// Flush all keys
func (s *Store) Flush() {
	s.store = make(map[string]*list.Element)
	s.ll = list.New()
}

// Stats returns up-to-date information about the cache
func (s *Store) Stats() string {
	// TODO (healeycodes)
	// Use json package here
	return fmt.Sprintf(`{"keyCount": %v, "maxSize": %v}`, len(s.store), s.max)
}
