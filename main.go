package main

import (
	"fmt"
	"sync"
	store "github.com/healeycodes/key-value-store/store"
)

func main() {
	keys := new(sync.Map)
	store.Set("Test", 1)
	fmt.Println("Test")
}
