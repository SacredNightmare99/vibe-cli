package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func HandleSave(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vibe save \"<your message>\"")
		return
	}
	message := args[0]

	patchData, err := exec.Command("git", "diff", "HEAD").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error running git diff: %v\n", err)
		os.Exit(1)
	}

	if len(patchData) == 0 {
		fmt.Println("[VIBE] No changes to save. Working directory is clean.")
		return
	}

	vibe := Vibe{
		ID:        time.Now().Unix(),
		Message:   message,
		Timestamp: time.Now().Format(time.RFC822),
	}

	patchFile := fmt.Sprintf("%s/%d.patch", PatchesDir, vibe.ID)
	_ = os.WriteFile(patchFile, patchData, 0644)

	vibes := ReadVibes()
	vibes = append(vibes, vibe)
	WriteVibes(vibes)

	fmt.Printf("[VIBE] 💾 Saved vibe [%d]: %s\n", vibe.ID, vibe.Message)
}

