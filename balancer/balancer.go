package balancer

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type Balancer struct {
	servers      []string
	current      int
	health       map[string]bool
	connections  map[string]int
	strategy     string
	requestPerIP map[string]int
}

func NewBalancer(servers []string) *Balancer {

	health := make(map[string]bool)
	connections := make(map[string]int)

	for _, server := range servers {
		health[server] = true
		connections[server] = 0
	}

	b := &Balancer{
		servers:      servers,
		current:      0,
		health:       health,
		connections:  connections,
		strategy:     "least-connections",
		requestPerIP: make(map[string]int),
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			b.requestPerIP = make(map[string]int)
		}
	}()

	return b

}

func (b *Balancer) SetStrategy(strategy string) {
	if strategy == "least-connections" || strategy == "round-robin" {
		b.strategy = strategy

		log.Printf("[STRATEGY SWITCH] strategy chenged to %s\n", strategy)
	} else {
		log.Printf("[STRATEGY SWITCH] unkown strategy %s\n", strategy)
	}
}

func (b *Balancer) getNextServer() string {

	start := b.current

	for {
		server := b.servers[b.current]
		b.current = (b.current + 1) % len(b.servers)
		if b.health[server] {
			return server
		}

		if b.current == start {
			return server
		}

	}
}

func (b *Balancer) GetLeastConenctionsServer() string {
	minConn := int(^uint(0) >> 1)
	var selected string

	for _, server := range b.servers {
		if b.health[server] && b.connections[server] <= minConn {
			minConn = b.connections[server]
			selected = server
		}
	}
	return selected
}

func (b *Balancer) ServeProxy(w http.ResponseWriter, r *http.Request) {

	ip := r.RemoteAddr
	ipOnly := ip

	if host, _, err := net.SplitHostPort(ip); err == nil {
		ipOnly = host
	}

	b.requestPerIP[ipOnly]++

	if b.requestPerIP[ipOnly] > 10 {
		log.Printf("[BLOCKED] IP: %s | Exceeded rate limit | URL: %s", ipOnly, r.RequestURI)
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}
	log.Printf("[REQUEST] IP: %s | URL: %s | Method: %s", ipOnly, r.RequestURI, r.Method)

	attempt := 0
	totalServers := len(b.servers)

	for attempt < totalServers {
		var targetServer string
		if b.strategy == "round-robin" {

			targetServer = b.getNextServer()
		} else {

			targetServer = b.GetLeastConenctionsServer()
		}
		b.connections[targetServer]++
		targetURL := targetServer + r.RequestURI

		req, err := http.NewRequest(r.Method, targetURL, r.Body)

		if err != nil {
			b.connections[targetServer]--
			http.Error(w, "Failed to Create Backend Request", http.StatusInternalServerError)
			return
		}

		req.Header = r.Header

		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		resp, err := client.Do(req)

		if err != nil || resp.StatusCode >= 500 {
			b.connections[targetServer]--
			log.Printf("[RETRY] IP: %s | Failed server: %s | Retrying another...", ipOnly, targetServer)
			attempt++

			if resp != nil {
				resp.Body.Close()
			}
			continue
		}

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		resp.Body.Close()
		log.Printf("[SUCCESS] IP: %s | Served by: %s | URL: %s | Status: %d", ipOnly, targetServer, r.RequestURI, resp.StatusCode)
		b.connections[targetServer]--
		return

	}
	log.Printf("[FAILURE] IP: %s | All servers down | URL: %s", ipOnly, r.RequestURI)
	http.Error(w, "All backend servers are down", http.StatusServiceUnavailable)
}

func (b *Balancer) StartHealthCheck(intervalSecs int) {
	go func() {
		for {
			for _, server := range b.servers {
				resp, err := http.Get(server + "/health")
				if err != nil || resp.StatusCode != http.StatusOK {
					b.health[server] = false
					log.Printf("[Health Check] server %s is DOWN\n", server)
				} else {
					b.health[server] = true
				}
				if resp != nil {
					resp.Body.Close()
				}
			}
			time.Sleep(time.Duration(intervalSecs) * time.Second)
		}
	}()
}
