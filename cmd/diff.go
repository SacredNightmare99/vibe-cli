package cmd

import (
	"fmt"
	"os"
)

func HandleDiff(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vibe diff <vibe_id>")
		return
	}
	vibeID := args[0]
	file := fmt.Sprintf("%s/%s.patch", PatchesDir, vibeID)

	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error reading patch: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

