package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "status":
		fmt.Println("sysix status - not yet implemented")
	case "watch":
		fmt.Println("sysix watch (TUI) - not yet implemented")
	case "serve":
		fmt.Println("sysix serve (web UI) - not yet implemented")
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
  sysix status     snapshot of your system right now
  sysix watch      launch the live TUI
  sysix serve      launch the web UI
  sysix help       show this help message
	`)
}