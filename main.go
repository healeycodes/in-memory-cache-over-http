package main

import (
	"net/http"
	"os"

	"github.com/healeycodes/key-value-store/store"
)

var s store.Store

func main() {
	s = store.StoreService()
	setup()
	start()
}

// Setup path handlers
func setup() {
	http.HandleFunc("/get", Get)
	http.HandleFunc("/set", Set)
}

// Listen on PORT. Defaults to 8000
func start() {
	err := http.ListenAndServe(getEnv("PORT", ":8000"), nil)
	if err != nil {
		panic(err)
	}
}

// Get a key from the store. Status code: 200 if present, else 404.
// e.g. ?key=foo
func Get(w http.ResponseWriter, r *http.Request) {
	value, exist := s.Get(r.URL.Query().Get("key"))
	if !exist {
		http.Error(w, "", 404)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(value.(string)))
}

// Set a key in the store. Status code: 204.
func Set(w http.ResponseWriter, r *http.Request) {
	s.Set(r.URL.Query().Get("key"), r.URL.Query().Get("value"))
	w.WriteHeader(http.StatusNoContent)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
