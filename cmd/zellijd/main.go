package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"lxc-terminal/internal/api"
	"lxc-terminal/internal/zellij"
)

func main() {
	var (
		port      = flag.Int("port", 8080, "Server port")
		container = flag.String("container", "", "Container name")
	)
	flag.Parse()

	// Initialize Zellij daemon
	daemon := zellij.NewDaemon(*container)

	// Initialize API handler
	apiHandler := api.NewHandler(daemon)

	// Setup routes
	mux := http.NewServeMux()

	// WebSocket endpoint for PTY
	mux.HandleFunc("/ws", apiHandler.HandleWebSocket)

	// REST API endpoints
	mux.HandleFunc("/api/sessions", apiHandler.HandleSessions)
	mux.HandleFunc("/api/sessions/", apiHandler.HandleSession)
	mux.HandleFunc("/api/health", apiHandler.HandleHealth)

	// Serve static files
	fs := http.FileServer(http.Dir("./web/"))
	mux.Handle("/", fs)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Zellijd starting on %s (container: %s)", addr, *container)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}