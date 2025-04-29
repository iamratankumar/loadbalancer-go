package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func startServer(port int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("Hello from the server on port %d", port)
		fmt.Fprint(w, message)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// if port == 5002 || port == 5003 || port == 5001 {
		// 	http.Error(w, "Server Down", http.StatusInternalServerError)
		// 	return
		// }
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	serverAddr := fmt.Sprintf(":%d", port)

	server := &http.Server{
		Addr:           serverAddr,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		log.Printf("Starting server on port %d...\n", port)

		if err := server.ListenAndServe(); err != nil {
			log.Printf("Server  on port %d failed %s\n", port, err)
		}
	}()
}
