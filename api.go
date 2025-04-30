package main

import (
	"bufio"
	"encoding/json"
	"loadbalancer-go/balancer"
	"net/http"
	"os"
	"path/filepath"
)

func StatusHandler(lb *balancer.Balancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type ServerStatus struct {
			Address     string `json: "address"`
			Health      bool   `json: "healthy"`
			Connections int    `json:"connections"`
		}

		var servers []ServerStatus
		for _, addr := range lb.Servers() {
			servers = append(servers, ServerStatus{
				Address:     addr,
				Health:      lb.IsHealthy(addr),
				Connections: lb.GetConnections(addr),
			})
		}
		response := map[string]interface{}{
			"strategy":      lb.GetStrategy(),
			"totalRequests": lb.GetTotalRequests(),
			"blockedIPs":    lb.GetBlockedIPs(),
			"servers":       servers,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func LogHandler(logPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logFile := filepath.Clean(logPath)
		file, err := os.Open(logFile)

		if err != nil {
			http.Error(w, "Could not open log file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		var lines []string

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		start := 0
		if len(lines) > 100 {
			start = len(lines) - 100
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lines[start:])

	}
}

func SetStrategyHandler(lb *balancer.Balancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Payload struct {
			Strategy string `json: "strategy"`
		}

		var data Payload

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		strategy := data.Strategy

		if strategy != "round-robin" && strategy != "least-connections" {
			http.Error(w, "Unsupported strategy", http.StatusBadRequest)
			return
		}

		lb.SetStrategy(strategy)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "success",
			"strategy": strategy,
		})
	}
}

func BlockedIpsHandler(lb *balancer.Balancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blocked := lb.GetBlockedIPs()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blocked)
	}
}
