package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "server":
		startServer()
	case "init":
		handleInit()
	case "save":
		handleSave(args)
	case "list":
		handleList()
	case "checkout":
		handleCheckout(args)
	case "diff":
		handleDiff(args)
	case "reset":
		handleReset()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`
Vibe - The AI-powered coding assistant and session manager.

Usage:
  vibe <command> [arguments]

Available Commands:
  server            Starts the API server for the mobile app.
  init              Initializes a vibegit session in the current project.
  save "<message>"  Manually saves your current changes as a new vibe.
  list              Lists all saved vibes in the current session.
  checkout <id>     Reverts the workspace to a previously saved vibe.
  diff <id>         Shows the changes captured in a specific vibe.
  reset             Clears the vibegit session and reverts all changes.
	`)
}
