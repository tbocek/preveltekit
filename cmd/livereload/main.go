package main

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	mu      sync.Mutex
	clients []chan struct{}
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch := make(chan struct{}, 1)
		mu.Lock()
		clients = append(clients, ch)
		mu.Unlock()

		defer func() {
			mu.Lock()
			for i, c := range clients {
				if c == ch {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			mu.Unlock()
		}()

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ch:
				fmt.Fprintf(w, "data: reload\n\n")
				f.Flush()
			}
		}
	})

	http.HandleFunc("/trigger", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		mu.Lock()
		for _, ch := range clients {
			select {
			case ch <- struct{}{}:
			default:
			}
		}
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Livereload server on :3001")
	http.ListenAndServe(":3001", nil)
}
