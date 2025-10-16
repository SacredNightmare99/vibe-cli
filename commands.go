// commands.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Vibe defines the structure for an entry in log.json
type Vibe struct {
	ID        int64  `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

const (
	vibesDir    = ".vibes"
	patchesDir  = ".vibes/patches"
	logFile     = ".vibes/log.json"
)

// handleInit sets up the .vibes directory and log file.
func handleInit() {
	fmt.Println("Initializing vibegit session...")
	if err := os.MkdirAll(patchesDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directories: %v\n", err)
		os.Exit(1)
	}

	// Create log file with an empty JSON array
	if err := os.WriteFile(logFile, []byte("[]"), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating log file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Vibegit initialized successfully!")
}

// handleSave captures current changes, saves a patch, and logs it.
func handleSave(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vibe save \"<your message>\"")
		return
	}
	message := args[0]

	// Run 'git diff HEAD' to get the changes
	patchData, err := exec.Command("git", "diff", "HEAD").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running git diff: %v\n", err)
		os.Exit(1)
	}

	if len(patchData) == 0 {
		fmt.Println("No changes to save. Working directory is clean.")
		return
	}

	// Create a new vibe
	newVibe := Vibe{
		ID:        time.Now().Unix(),
		Message:   message,
		Timestamp: time.Now().Format(time.RFC822),
	}

	// Save the patch file
	patchFileName := fmt.Sprintf("%s/%d.patch", patchesDir, newVibe.ID)
	if err := os.WriteFile(patchFileName, patchData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing patch file: %v\n", err)
		os.Exit(1)
	}

	// Read existing vibes, add the new one, and write back
	vibes := readVibes()
	vibes = append(vibes, newVibe)
	writeVibes(vibes)

	fmt.Printf("💾 Vibe [%d] saved: %s\n", newVibe.ID, newVibe.Message)
}

// handleList displays all saved vibes.
func handleList() {
	vibes := readVibes()
	if len(vibes) == 0 {
		fmt.Println("No vibes saved yet. Use 'vibe save \"<message>\"' to save one.")
		return
	}

	fmt.Println("--- VIBE HISTORY ---")
	// Print newest first
	for i := len(vibes) - 1; i >= 0; i-- {
		vibe := vibes[i]
		fmt.Printf("ID: %d | %s | %s\n", vibe.ID, vibe.Timestamp, vibe.Message)
	}
	fmt.Println("--------------------")
}

// handleCheckout reverts the working directory to a saved vibe.
func handleCheckout(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vibe checkout <vibe_id>")
		return
	}
	vibeID := args[0]
	patchFile := fmt.Sprintf("%s/%s.patch", patchesDir, vibeID)

	// Check if the patch file actually exists
	if _, err := os.Stat(patchFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Vibe with ID %s not found.\n", vibeID)
		os.Exit(1)
	}

	// 1. Clean the working directory
	if err := exec.Command("git", "checkout", "--", ".").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error cleaning workspace: %v\n", err)
		os.Exit(1)
	}

	// 2. Apply the patch
	// Using -3way to handle potential conflicts more gracefully
	cmd := exec.Command("git", "apply", "--3way", patchFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "Error applying patch:\n%s\n", output)
		os.Exit(1)
	}

	fmt.Printf("⏪ Workspace reverted to vibe [%s].\n", vibeID)
}

// handleDiff shows the contents of a specific patch file.
func handleDiff(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vibe diff <vibe_id>")
		return
	}
	vibeID := args[0]
	patchFile := fmt.Sprintf("%s/%s.patch", patchesDir, vibeID)

	patchData, err := os.ReadFile(patchFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not read vibe with ID %s.\n", vibeID)
		os.Exit(1)
	}

	fmt.Println(string(patchData))
}

// handleReset cleans the workspace and removes the .vibes directory.
func handleReset() {
	fmt.Println("🧹 Cleaning workspace and removing session...")

	// 1. Clean the working directory
	if err := exec.Command("git", "checkout", "--", ".").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error cleaning workspace: %v\n", err)
		os.Exit(1)
	}

	// 2. Remove the .vibes directory
	if err := os.RemoveAll(vibesDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing .vibes directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("💥 Vibegit session cleared.")
}

// --- HELPER FUNCTIONS ---

// readVibes is a helper to read and parse the log.json file.
func readVibes() []Vibe {
	file, err := os.ReadFile(logFile)
	if err != nil {
		// If the file doesn't exist, it means we're not initialized.
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Error: vibegit not initialized. Run 'vibe init' first.")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error reading log file: %v\n", err)
		os.Exit(1)
	}

	var vibes []Vibe
	if err := json.Unmarshal(file, &vibes); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing log file: %v\n", err)
		os.Exit(1)
	}
	return vibes
}

// writeVibes is a helper to serialize and write vibes to the log.json file.
func writeVibes(vibes []Vibe) {
	// Marshal with indentation for human readability
	data, err := json.MarshalIndent(vibes, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling vibes: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(logFile, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
		os.Exit(1)
	}
}
