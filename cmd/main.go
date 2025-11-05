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
	case "watch":
		handleWatch()
	default:
		fmt.Printf("[VIBE] Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`
Vibe - Local AI coding companion

Usage:
  vibe <command> [arguments]

Available Commands:
  init              Initialize and mark current project
  save "<msg>"      Save current git diff as vibe patch
  list              List saved vibes
  checkout <id>     Revert workspace to vibe
  diff <id>         Show changes in vibe
  reset             Clear session and patches
  watch             Start watcher daemon
`)
}

