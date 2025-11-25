package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var counter int
var lock sync.Mutex
var leak = []string{}

func handler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	counter++
	lock.Unlock()

	time.Sleep(2 * time.Second)

	leak = append(leak, fmt.Sprintf("req-%d", counter))

	fmt.Fprintf(w, "counter=%d", counter)
}

func health(w http.ResponseWriter, r *http.Request) {
	if time.Now().Unix()%2 == 0 {
		w.WriteHeader(500)
		return
	}
	w.Write([]byte("ok"))
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", health)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

