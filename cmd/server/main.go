package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Serve static files from web directory
	fs := http.FileServer(http.Dir("./web/"))
	http.Handle("/", fs)

	// WebSocket endpoint (placeholder for T2.1)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "WebSocket endpoint - not implemented yet (T2.1)")
	})

	fmt.Println("LXC Terminal PWA Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}