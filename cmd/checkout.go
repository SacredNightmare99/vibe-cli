package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

func HandleCheckout(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vibe checkout <vibe_id>")
		return
	}
	id := args[0]
	file := fmt.Sprintf("%s/%s.patch", PatchesDir, id)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "[VIBE] Patch not found: %s\n", id)
		os.Exit(1)
	}

	_ = exec.Command("git", "checkout", "--", ".").Run()
	cmd := exec.Command("git", "apply", "--3way", file)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error applying patch:\n%s\n", output)
		os.Exit(1)
	}
	fmt.Printf("[VIBE] ⏪ Workspace reverted to vibe [%s].\n", id)
}

