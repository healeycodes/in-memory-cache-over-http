package main

import (
	"fmt"
	"sync"
	store "github.com/healeycodes/key-value-store/store"
)

keys := new(sync.Map)
func main() {
	store.Set("Test", 1)
	fmt.Println("Test")
}