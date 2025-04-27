package main

import (
	"fmt"
	"net/http"
)

func startServer(port int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("Hello from the server on port %d", port)
		fmt.Fprint(w, message)
	})

	serverAddr := fmt.Sprintf(":%d", port)

	go func() {
		fmt.Printf("Starting server on port %d...\n", port)
		if err := http.ListenAndServe(serverAddr, mux); err != nil {
			fmt.Printf("Server  on port %d failed %s\n", port, err)
		}
	}()
}
