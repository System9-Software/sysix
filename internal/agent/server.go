package agent

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/System9-Software/sysix/internal/collector"
)

func Start(host string, port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/snapshot", handleSnapshot)
	mux.HandleFunc("/api/ports", handlePorts)
	mux.HandleFunc("/api/network", handleNetwork)

	if host == "" {
		host = "0.0.0.0"
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("sysix agent running at http://%s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("sysix agent"))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleSnapshot(w http.ResponseWriter, r *http.Request) {
	snapshot, err := collector.GetSnapshot()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(snapshot)
}

func handlePorts(w http.ResponseWriter, r *http.Request) {
	ports, err := collector.GetPorts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ports)
}

func handleNetwork(w http.ResponseWriter, r *http.Request) {
	network, err := collector.GetNetwork()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(network)
}
