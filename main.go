package main

import (
	"fmt"
	"loadbalancer-go/balancer"
	"log"
	"net/http"
	"time"

	"github.com/natefinch/lumberjack"
)

func main() {

	log.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/loadbalancer.log",
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 100,
		Compress:   false,
	})

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[STARTUP] Load Balancer is starting...")

	// Force a manual rotation every 10 minutes
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			log.Println("[SNAPSHOT] 10 minutes passed, rotating logs manually...")
			log.SetOutput(&lumberjack.Logger{
				Filename:   "./logs/loadbalancer.log",
				MaxSize:    10,
				MaxAge:     1,
				MaxBackups: 100,
				Compress:   false,
			})
		}
	}()

	startServer(5001)
	startServer(5002)
	startServer(5003)

	servers := []string{
		"http://localhost:5001",
		"http://localhost:5002",
		"http://localhost:5003",
	}

	lb := balancer.NewBalancer(servers)
	lb.StartHealthCheck(5)
	lb.SetStrategy("round-robin")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(w, r)
	})

	fmt.Println("Load balancer started on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("Load balancer failed %s\n", err)
	}

}
