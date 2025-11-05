package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

func HandleReset() {
	fmt.Println("[VIBE] 🧹 Cleaning workspace and removing session...")

	_ = exec.Command("git", "checkout", "--", ".").Run()
	_ = os.RemoveAll(VibesDir)

	fmt.Println("[VIBE] 💥 Vibegit session cleared.")
}

