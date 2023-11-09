package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// Handler function for the '/' route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Only respond to GET requests
		if r.Method == http.MethodGet {
			// Call the upstream server
			resp, err := http.Get("http://server-service.default.svc.cluster.local:5000")
			if err != nil {
				// Handle error.
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				// Handle error.
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Write the response body to the response writer
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		} else {
			// If the request method is not GET, return a 405 Method Not Allowed
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Log that the server is starting
	fmt.Println("Server is listening on port 6000...")

	// Listen on port 6000
	if err := http.ListenAndServe(":6000", nil); err != nil {
		log.Fatal(err)
	}
}
