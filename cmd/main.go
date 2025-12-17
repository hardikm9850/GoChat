package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Staring GoChat")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Program blocks here and waits for HTTP requests
	log.Fatal(http.ListenAndServe(":8080", nil))
}
