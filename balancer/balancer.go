package balancer

import (
	"io"
	"net/http"
)

type Balancer struct {
	servers []string
	current int
}

func NewBalancer(servers []string) *Balancer {
	return &Balancer{
		servers: servers,
		current: 0,
	}
}

func (b *Balancer) getNextServer() string {

	server := b.servers[b.current]
	b.current = (b.current + 1) % len(b.servers)
	return server
}

func (b *Balancer) ServeProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := b.getNextServer()
	targetUrl := targetServer + r.RequestURI

	req, err := http.NewRequest(r.Method, targetUrl, r.Body)

	if err != nil {
		http.Error(w, "failed to create request to backend", http.StatusInternalServerError)
	}

	req.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		http.Error(w, "Failed to reach to backend", http.StatusServiceUnavailable)
	}

	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}
