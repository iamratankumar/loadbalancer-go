package main

import (
	"fmt"
	"loadbalancer-go/balancer"
	"net/http"
)

func main() {
	startServer(5001)
	startServer(5002)
	startServer(5003)

	servers := []string{
		"http://localhost:5001",
		"http://localhost:5002",
		"http://localhost:5003",
	}

	lb := balancer.NewBalancer(servers)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(w, r)
	})

	fmt.Println("Load balancer started on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Load balancer failed %s\n", err)
	}

}
