package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "pong"}); err != nil {
			log.Printf("Failed to encode ping response: %v", err)
		}
	})

	log.Println("Test API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
