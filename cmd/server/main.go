package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/pz7-redis/internal/cache"
	"github.com/redis/go-redis/v9"
)

func main() {
	c := cache.New("localhost:6379")
	defer c.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		key, value := r.URL.Query().Get("key"), r.URL.Query().Get("value")
		if key == "" || value == "" {
			http.Error(w, "key and value required", http.StatusBadRequest)
			return
		}
		if err := c.Set(key, value, 10*time.Second); err != nil {
			http.Error(w, "redis unavailable", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "OK: %s=%s (TTL 10s)", key, value)
	})
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key required", http.StatusBadRequest)
			return
		}
		value, err := c.Get(key)
		if errors.Is(err, redis.Nil) {
			http.Error(w, "key not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "redis unavailable", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "VALUE: %s=%s", key, value)
	})
	mux.HandleFunc("/ttl", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key required", http.StatusBadRequest)
			return
		}
		ttl, err := c.TTL(key)
		if err != nil {
			http.Error(w, "redis unavailable", http.StatusInternalServerError)
			return
		}
		if ttl < 0 {
			fmt.Fprintf(w, "TTL for %s: key not found (-2)", key)
			return
		}
		fmt.Fprintf(w, "TTL for %s: %v", key, ttl)
	})
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
