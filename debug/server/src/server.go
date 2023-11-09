package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've reached the server!")
	})

	fmt.Println("Server is listening on port 5000...")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}
