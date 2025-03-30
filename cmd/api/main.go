package main

import (
	"empshift-csp/internal/api"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Println("Server running on port 8080")
	log.Fatal(server.ListenAndServe())
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/schedule", api.HandleScheduleRequest)
	mux.HandleFunc("/", handleNotFound)
	return mux
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Welcome to Scheduling API"))
}
