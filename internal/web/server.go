package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/System9-Software/sysix/internal/collector"
)

func Start(port int) error {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/snapshot", handleSnapshot)
	http.HandleFunc("/api/ports", handlePorts)
	http.HandleFunc("/api/network", handleNetwork)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("sysix web UI running at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>sysix</h1><p>web UI coming soon</p>"))
}

func handleSnapshot(w http.ResponseWriter, r *http.Request) {
	snapshot, err := collector.GetSnapshot()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshot)
}

func handlePorts(w http.ResponseWriter, r *http.Request) {
	ports, err := collector.GetPorts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ports)
}

func handleNetwork(w http.ResponseWriter, r *http.Request) {
	network, err := collector.GetNetwork()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(network)
}
