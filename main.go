package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/healeycodes/key-value-store/store"
)

var s *store.Store

func main() {
	New()
	Setup()
	Start()
}

// Setup path handlers
func Setup() {
	http.HandleFunc("/get", Get)
	http.HandleFunc("/set", Set)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/checkandset", CheckAndSet)
	http.HandleFunc("/increment", Increment)
	http.HandleFunc("/decrement", Decrement)
	http.HandleFunc("/append", Append)
	http.HandleFunc("/prepend", Prepend)
	http.HandleFunc("/flush", Flush)
}

// Start Listening on PORT. Defaults to 8000.
func Start() {
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

// CheckAndSet a key in the store if it matches the compare value.
// Status code: 204 if matches else 400
func CheckAndSet(w http.ResponseWriter, r *http.Request) {
	s.CheckAndSet(r.URL.Query().Get("key"), r.URL.Query().Get("value"), r.URL.Query().Get("compare"))
	w.WriteHeader(http.StatusNoContent)
}

// Increment a key in the store by an amount. If key missing, set the amount.
// Status code: 204 if incrementable else 400.
func Increment(w http.ResponseWriter, r *http.Request) {
	if err := s.Increment(r.URL.Query().Get("key"), r.URL.Query().Get("value")); err == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Decrement a key in the store by an amount. If key missing, set the amount.
// Status code: 204 if decrementable else 400.
func Decrement(w http.ResponseWriter, r *http.Request) {
	if err := s.Decrement(r.URL.Query().Get("key"), r.URL.Query().Get("value")); err == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Append to a key in the store.
// Status code: 204 if key exists else 400.
func Append(w http.ResponseWriter, r *http.Request) {
	if exists := s.Append(r.URL.Query().Get("key"), r.URL.Query().Get("value")); exists == true {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Prepend to a key in the store.
// Status code: 204 if key exists else 400.
func Prepend(w http.ResponseWriter, r *http.Request) {
	if exists := s.Prepend(r.URL.Query().Get("key"), r.URL.Query().Get("value")); exists == true {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Flush all keys.
func Flush(w http.ResponseWriter, r *http.Request) {
	s.Flush()
	w.WriteHeader(http.StatusNoContent)
}

// New store.
func New() {
	size, _ := strconv.Atoi(getEnv("SIZE", "2"))
	s = store.Service(size)
}

// Gets an ENV variable else returns fallback.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
