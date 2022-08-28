package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	uploadsDir = STATIC_DIR + "/uploads"
)

func runServer(srv_host string) {
	// Define endpoints
	// FileServer
	fileServer := http.FileServer(http.Dir(STATIC_DIR))
	http.Handle("/", fileServer)
	http.HandleFunc("/health", handleHealthCheck)

	// Start the server
	log.Printf("Starting web server on: %s\n", srv_host)
	if err := http.ListenAndServe(srv_host, nil); err != nil {
		log.Fatal("Could not start server: ", err)
	}
}

// Health check
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Printf("Request method recieved was not GET\n")
		log.Printf("Request method: %s\n", r.Method)
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}

	log.Printf("Got request for /health - responding with Healthy")
	fmt.Fprintf(w, "Healthy")
}
