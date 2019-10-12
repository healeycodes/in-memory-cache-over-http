package main

import (
	"net/http"
	"os"

	"github.com/healeycodes/key-value-store/store"
)

var s store.Store

func main() {
	s = store.Service()
	setup()
	start()
}

// Setup path handlers
func setup() {
	http.HandleFunc("/get", Get)
	http.HandleFunc("/set", Set)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/increment", Increment)
	http.HandleFunc("/decrement", Decrement)
	http.HandleFunc("/flush", Flush)
}

// Listen on PORT. Defaults to 8000.
func start() {
	err := http.ListenAndServe(getEnv("PORT", ":8000"), nil)
	if err != nil {
		panic(err)
	}
}

// Get a key from the store.
// Status code: 200 if present, else 404.
// e.g. ?key=foo
func Get(w http.ResponseWriter, r *http.Request) {
	value, exist := s.Get(r.URL.Query().Get("key"))
	if !exist {
		http.Error(w, "", 404)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(value))
}

// Set a key in the store.
// Status code: 204.
func Set(w http.ResponseWriter, r *http.Request) {
	s.Set(r.URL.Query().Get("key"), r.URL.Query().Get("value"))
	w.WriteHeader(http.StatusNoContent)
}

// Delete a key in the store.
// Status code: 204.
func Delete(w http.ResponseWriter, r *http.Request) {
	s.Delete(r.URL.Query().Get("key"))
	w.WriteHeader(http.StatusNoContent)
}

// Increment a key in the store. Sets 1 if new key.
// Status code: 204 if incrementable else 400.
func Increment(w http.ResponseWriter, r *http.Request) {
	if err := s.Increment(r.URL.Query().Get("key")); err == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Decrement a key in the store. Sets -1 if new key.
// Status code: 204 if decrementable else 400.
func Decrement(w http.ResponseWriter, r *http.Request) {
	if err := s.Decrement(r.URL.Query().Get("key")); err == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Append to a key in the store.
// Status code: 204 if key exists else 400.
func Append(w http.ResponseWriter, r *http.Request) {
	if exists := s.Append(r.URL.Query().Get("key"), r.URL.Query().Get("value")); exists == false {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Flush all keys in the store.
// Status code: 204 if key exists else 400.
func Flush(w http.ResponseWriter, r *http.Request) {
	s = store.Service()
}

// Gets an ENV variable else returns fallback.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
