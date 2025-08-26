package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Create a new ServeMux to handle requests
	mux := http.NewServeMux()

	// Handle all requests with a logging middleware
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log request details
		log.Printf(
			"[%s] %s %s %s from %s (User-Agent: %s)",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			r.Proto,
			r.RemoteAddr,
			r.UserAgent(),
		)

		// Simple response
		fmt.Fprintf(w, "Request received and logged")
	})

	// Create server
	server := &http.Server{
		Addr:    ":33333",
		Handler: mux,
	}

	// Start server and log any errors
	log.Printf("Starting server on port 33333...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
