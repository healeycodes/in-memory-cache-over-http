package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"healeycodes/in-memory-cache-over-http/cache"
)

var s *cache.Store

// Listen on PORT. Defaults to 8000
func Listen() {
	new()
	setup()
	start()
}

// New store
func new() {
	size, _ := strconv.Atoi(getEnv("SIZE", "0"))
	s = cache.Service(size)
}

// Setup path handlers
func setup() {
	http.HandleFunc("/get", handle(Get))
	http.HandleFunc("/set", handle(Set))
	http.HandleFunc("/delete", handle(Delete))
	http.HandleFunc("/checkandset", handle(CheckAndSet))
	http.HandleFunc("/increment", handle(Increment))
	http.HandleFunc("/decrement", handle(Decrement))
	http.HandleFunc("/append", handle(Append))
	http.HandleFunc("/prepend", handle(Prepend))
	http.HandleFunc("/stats", handle(Stats))
	http.HandleFunc("/flush", handle(Flush))
}

// Start http
func start() {
	port := getEnv("PORT", ":8000")
	fmt.Println("Listening on", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}

// Get a key from the store
// Status code: 200 if present, else 404
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

// Set a key in the store
// Status code: 204
func Set(w http.ResponseWriter, r *http.Request) {
	s.Set(
		r.URL.Query().Get("key"),
		r.URL.Query().Get("value"),
		getExpire(r.URL.Query().Get("expire")))
	w.WriteHeader(http.StatusNoContent)
}

// Delete a key in the store
// Status code: 204
func Delete(w http.ResponseWriter, r *http.Request) {
	s.Delete(r.URL.Query().Get("key"))
	w.WriteHeader(http.StatusNoContent)
}

// CheckAndSet a key in the store if it matches the compare value
// Status code: 204 if matches else 400
func CheckAndSet(w http.ResponseWriter, r *http.Request) {
	if s.CheckAndSet(
		r.URL.Query().Get("key"),
		r.URL.Query().Get("value"),
		getExpire(r.URL.Query().Get("expire")),
		r.URL.Query().Get("compare")) == true {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// Increment a key in the store by an amount. If key missing, set the amount
// Status code: 204 if incrementable else 400
func Increment(w http.ResponseWriter, r *http.Request) {
	if err := s.Increment(
		r.URL.Query().Get("key"),
		r.URL.Query().Get("value"),
		getExpire(r.URL.Query().Get("expire"))); err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// Decrement a key in the store by an amount. If key missing, set the amount
// Status code: 204 if decrementable else 400
func Decrement(w http.ResponseWriter, r *http.Request) {
	if err := s.Decrement(
		r.URL.Query().Get("key"),
		r.URL.Query().Get("value"),
		getExpire(r.URL.Query().Get("expire"))); err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// Append to a key in the store
// Status code: 204
func Append(w http.ResponseWriter, r *http.Request) {
	s.Append(
		r.URL.Query().Get("key"),
		r.URL.Query().Get("value"),
		getExpire(r.URL.Query().Get("expire")))
	w.WriteHeader(http.StatusNoContent)
}

// Prepend to a key in the store
// Status code: 204
func Prepend(w http.ResponseWriter, r *http.Request) {
	s.Prepend(
		r.URL.Query().Get("key"),
		r.URL.Query().Get("value"),
		getExpire(r.URL.Query().Get("expire")))
}

// Flush all keys
func Flush(w http.ResponseWriter, r *http.Request) {
	s.Flush()
	w.WriteHeader(http.StatusNoContent)
}

// Stats of the cache
// Status code: 200
func Stats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(s.Stats()))
}

// Middleware
func handle(f func(http.ResponseWriter, *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Mutex.Lock()
		defer s.Mutex.Unlock()
		if getEnv("APP_ENV", "") != "production" {
			fmt.Println(time.Now(), r.URL)
		}
		f(w, r)
	}
}

// Safely get the expire, 0 if error
func getExpire(attempt string) int {
	value, err := strconv.Atoi(attempt)
	if err != nil {
		return 0
	}
	return value
}

// Gets an ENV variable else returns fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
