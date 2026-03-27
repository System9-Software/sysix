package main

import (
	"fmt"
	"os"
	"time"

	"github.com/System9-Software/sysix/internal/agent"
	"github.com/System9-Software/sysix/internal/collector"
	"github.com/System9-Software/sysix/internal/config"
	"github.com/System9-Software/sysix/internal/tui"
	"github.com/System9-Software/sysix/internal/web"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "status":
		showProcs := false
		showPorts := false
		for _, arg := range args[1:] {
			if arg == "--procs" {
				showProcs = true
			}
			if arg == "--ports" {
				showPorts = true
			}
		}

		snapshot, err := collector.GetSnapshot()
		if err != nil {
			fmt.Println("error collecting system data:", err)
			return
		}
		fmt.Printf("\n[*] Host:   %s (%s)\n", snapshot.Hostname, snapshot.OS)
		fmt.Printf("[~] Uptime: %d hours\n", snapshot.Uptime/3600)
		fmt.Printf("[>] CPU:    %.1f%%\n", snapshot.CPUPercent)
		fmt.Printf("[>] Memory: %.1f%% (%d MB / %d MB)\n", snapshot.MemPercent, snapshot.MemUsed/1024/1024, snapshot.MemTotal/1024/1024)
		fmt.Printf("[>] Disk:   %.1f%% (%d GB / %d GB)\n", snapshot.DiskPercent, snapshot.DiskUsed/1024/1024/1024, snapshot.DiskTotal/1024/1024/1024)

		if showProcs {
			fmt.Println("\n--- Processes ---")
			fmt.Printf("%-6s %-30s %8s %10s %s\n", "PID", "Name", "CPU%", "Mem(MB)", "Status")
			fmt.Println("----------------------------------------------------------------------")
			for _, p := range snapshot.Processes {
				if p.MemMB > 1 {
					fmt.Printf("%-6d %-30s %8.1f %10.1f %s\n", p.PID, p.Name, p.CPUPercent, p.MemMB, p.Status)
				}
			}
		}

		if showPorts {
			ports, err := collector.GetPorts()
			if err != nil {
				fmt.Println("error collecting ports:", err)
			} else {
				fmt.Println("\n--- Open Ports ---")
				fmt.Printf("%-8s %-6s %-8s %s\n", "Port", "Type", "Status", "PID")
				fmt.Println("--------------------------------")
				for _, p := range ports {
					fmt.Printf("%-8d %-6s %-8s %d\n", p.Port, p.Type, p.Status, p.PID)
				}
			}
		}

	case "watch":
		showProcs := false
		showPorts := false
		for _, arg := range args[1:] {
			if arg == "--procs" {
				showProcs = true
			}
			if arg == "--ports" {
				showPorts = true
			}
		}
		if showProcs || showPorts {
			if err := tui.StartFocused(showProcs, showPorts); err != nil {
				fmt.Println("error starting focused TUI:", err)
			}
		} else {
			if err := tui.Start(); err != nil {
				fmt.Println("error starting TUI:", err)
			}
		}
	case "serve":
		cfg, _ := config.Load()
		if !cfg.Web.Enabled {
			fmt.Println("web UI is disabled in config.yaml")
			return
		}
		if err := web.Start(cfg.Web.Port); err != nil {
			fmt.Println("error starting web server:", err)
		}
	case "agent":
		cfg, _ := config.Load()
		if !cfg.Agent.Enabled {
			fmt.Println("agent is disabled in config.yaml")
			return
		}
		if err := agent.Start(cfg.Agent.Host, cfg.Agent.Port); err != nil {
			fmt.Println("error starting agent server:", err)
		}
	case "observer":
		cfg, _ := config.Load()
		if !cfg.Web.Enabled {
			fmt.Println("web UI is disabled in config.yaml")
			return
		}
		targets := make([]web.AgentTarget, 0, len(cfg.Observer.Agents))
		for i, a := range cfg.Observer.Agents {
			if !a.Enabled || a.URL == "" {
				continue
			}
			id := a.ID
			if id == "" {
				id = fmt.Sprintf("agent-%d", i+1)
			}
			name := a.Name
			if name == "" {
				name = id
			}
			targets = append(targets, web.AgentTarget{
				ID:   id,
				Name: name,
				URL:  a.URL,
			})
		}
		poll := time.Duration(cfg.Observer.PollInterval) * time.Second
		if poll <= 0 {
			poll = 2 * time.Second
		}
		if err := web.StartObserver(cfg.Web.Port, targets, poll); err != nil {
			fmt.Println("error starting observer web server:", err)
		}
	case "help":
		printHelp()
	default:
		fmt.Printf("unknown command: %s\n", args[0])
		printHelp()
	}
}

func printHelp() {
	fmt.Println(`
welcome to sysix observer 0.3

Usage:
  sysix status            snapshot of your system right now
  sysix status --procs    launch focused process viewer
  sysix status --ports    launch focused port viewer
  sysix status --procs --ports  launch side-by-side viewer
  sysix watch             launch the live TUI
  sysix agent             run sysix agent API server
  sysix observer          run sysix observer (web + multi-host polling)
  sysix serve             launch the web UI
  sysix help              show this help message
	`)
}
