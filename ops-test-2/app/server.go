package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	_ "net/http/pprof" // pprof
)

type safeCounter struct {
	mu sync.RWMutex
	m  map[string]int
}

func newSafeCounter() *safeCounter {
	return &safeCounter{m: make(map[string]int)}
}

func (s *safeCounter) Inc(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[k]++
}
func (s *safeCounter) Get(k string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m[k]
}

var requests uint64
var counter = newSafeCounter()

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmsgprefix)
	mux := http.NewServeMux()

	// health endpoint - deterministic
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// simple health: server always healthy unless shutting down
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// simple metrics
		c := atomic.LoadUint64(&requests)
		resp := map[string]uint64{"requests": c}
		_ = json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&requests, 1)
		// simulate bounded async work; use context to cancel work when request cancelled
		ctx := r.Context()

		// Do not spawn long-running goroutines per-request without cancelation:
		errChan := make(chan error, 1)
		go func() {
			// simulate work; short sleep to not block
			select {
			case <-time.After(100 * time.Millisecond): // adjust latency
				counter.Inc("work")
				errChan <- nil
			case <-ctx.Done():
				errChan <- ctx.Err()
			}
		}()

		select {
		case err := <-errChan:
			if err != nil {
				logger.Printf("request cancelled: %v", err)
				http.Error(w, "request cancelled", http.StatusRequestTimeout)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("done"))
		case <-time.After(5 * time.Second):
			// bounded timeout to avoid hanging
			logger.Println("request timeout")
			http.Error(w, "timeout", http.StatusGatewayTimeout)
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(logger, mux),
		// Read/Write timeouts to prevent slow clients
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Println("starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server failed: %v", err)
		}
	}()

	// pprof listener (optional) - helpful during debugging
	go func() {
		logger.Println("pprof on :6060")
		log.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logger.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	logger.Println("shutdown complete")
}

func loggingMiddleware(logger *log.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		h.ServeHTTP(w, r)
		logger.Printf("completed in %s", time.Since(start))
	})
}
