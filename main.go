package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Rahul207158/Bitcask_GO/api"
	"github.com/Rahul207158/Bitcask_GO/kvstore"
)

func main() {
	// Create data directory if it doesn't exist
	dataDir := "data"

	// Initialize the store
	store, err := kvstore.NewStore(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}

	// Create API server
	server := api.NewServer(store)

	// Set up HTTP handlers
	http.HandleFunc("/put", server.PutHandler)
	http.HandleFunc("/get", server.GetHandler)

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nShutting down...")
		if err := store.Close(); err != nil {
			log.Printf("Error closing store: %v", err)
		}
		os.Exit(0)
	}()

	// Start the server
	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
