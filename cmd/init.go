package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

func HandleInit() {
	fmt.Println("[VIBE] Initializing vibegit session...")

	if err := os.MkdirAll(PatchesDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error creating directories: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(LogFile); os.IsNotExist(err) {
		_ = os.WriteFile(LogFile, []byte("[]"), 0644)
	}

	cwd, _ := os.Getwd()
	var tracked Tracked

	if data, err := os.ReadFile(TrackedFile); err == nil {
		_ = json.Unmarshal(data, &tracked)
	}
	for _, dir := range tracked.Projects {
		if dir == cwd {
			fmt.Println("[VIBE] 📍 Already tracked:", cwd)
			return
		}
	}
	tracked.Projects = append(tracked.Projects, cwd)
	data, _ := json.MarshalIndent(tracked, "", "  ")
	_ = os.WriteFile(TrackedFile, data, 0644)

	fmt.Println("[VIBE] 📍 Marked current directory:", cwd)
	fmt.Println("[VIBE] ✅ Vibegit initialized successfully.")
}

