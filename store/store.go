package github.com/healeycodes/key-value-store/store

import (
	"sync"
)

var keys sync.Map

type values struct {
	data []interface{}
}

// Set a key
func Set(key string, value interface{}) {
	keys.Store(key, value)
}

// Get a key
func Get(key string) (value interface{}, ok bool) {
	return keys.Load(key)
}

// Delete a key
func Delete(key string) {
	keys.Delete(key)
}

// Returns the values of all specified keys
// nil is used where the key doesn't exist
func Mget(keys []string) values {
	tmp := values{}
	for i := range keys {
		val, has := Get(keys[i])
		if !has {
			tmp.data = append(tmp.data, nil)
		} else {
			tmp.data = append(tmp.data, val)
		}
	}
	return tmp
}

// Sets the values of all specified keys
// array is parsed as <key> <value> repeated
func Mset(keysValues []string) {
	for i := range keysValues {
		if i%2 == 0 {
			Set(keysValues[i], keysValues[i+1])
		}
	}
}

// Clear everything
func Flush() {
	keys := new(sync.Map)
}
