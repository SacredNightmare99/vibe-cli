package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func HandleInit(args []string) {
	fmt.Println("[VIBE] Initializing vibegit session...")

	// 1. Create local .vibes structure in CWD
	if err := os.MkdirAll(PatchesDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error creating local directories: %v\n", err)
		os.Exit(1)
	}
	if _, err := os.Stat(LogFile); os.IsNotExist(err) {
		_ = os.WriteFile(LogFile, []byte("[]"), 0644)
	}

	// 2. Add project to central tracking file
	cwd, _ := os.Getwd()
	var projectID string
	if len(args) > 0 {
		projectID = args[0]
	} else {
		projectID = filepath.Base(cwd) 
	}

	tracked, err := ReadTrackedProjects()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error reading central tracked projects: %v\n", err)
		os.Exit(1)
	}

	// Check for duplicates
	for _, p := range tracked.Projects {
		if p.Path == cwd {
			fmt.Println("[VIBE] 📍 Project already tracked:", p.Path)
			return
		}
		if p.ID == projectID {
			fmt.Fprintf(os.Stderr, "[VIBE] Error: Project ID '%s' already exists for path %s\n", projectID, p.Path)
			os.Exit(1)
		}
	}

	// Add new project
	newProject := Project{
		ID:   projectID,
		Path: cwd,
	}
	tracked.Projects = append(tracked.Projects, newProject)
	if err := WriteTrackedProjects(tracked); err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error writing central tracked projects: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("[VIBE] 📍 Marked current directory as project:", projectID)
	fmt.Println("[VIBE] ✅ Vibegit initialized successfully.")
}
