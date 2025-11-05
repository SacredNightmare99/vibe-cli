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
	trackedFile = ".vibes/tracked.json"
)

// handleInit sets up the .vibes directory, log file, and marks current project.
func handleInit() {
	fmt.Println("[VIBE] Initializing vibegit session...")

	// Ensure .vibes structure exists
	if err := os.MkdirAll(patchesDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error creating directories: %v\n", err)
		os.Exit(1)
	}

	// Create empty log file if missing
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		if err := os.WriteFile(logFile, []byte("[]"), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "[VIBE] Error creating log file: %v\n", err)
			os.Exit(1)
		}
	}

	// --- Automatically mark current directory ---
	cwd, _ := os.Getwd()

	var tracked struct {
		Projects []string `json:"projects"`
	}

	// Load existing tracked file if present
	if data, err := os.ReadFile(trackedFile); err == nil {
		_ = json.Unmarshal(data, &tracked)
	}

	// Check if current directory already tracked
	for _, dir := range tracked.Projects {
		if dir == cwd {
			fmt.Println("[VIBE] 📍 Current directory already tracked:", cwd)
			fmt.Println("[VIBE] ✅ Vibegit initialized successfully!")
			return
		}
	}

	// Append and save
	tracked.Projects = append(tracked.Projects, cwd)
	data, _ := json.MarshalIndent(tracked, "", "  ")
	_ = os.WriteFile(trackedFile, data, 0644)

	fmt.Println("[VIBE] 📍 Marked current directory as tracked project:", cwd)
	fmt.Println("[VIBE] ✅ Vibegit initialized successfully!")
}

// handleSave captures current changes, saves a patch, and logs it.
func handleSave(args []string) {
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

	newVibe := Vibe{
		ID:        time.Now().Unix(),
		Message:   message,
		Timestamp: time.Now().Format(time.RFC822),
	}

	patchFileName := fmt.Sprintf("%s/%d.patch", patchesDir, newVibe.ID)
	if err := os.WriteFile(patchFileName, patchData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error writing patch file: %v\n", err)
		os.Exit(1)
	}

	vibes := readVibes()
	vibes = append(vibes, newVibe)
	writeVibes(vibes)

	fmt.Printf("[VIBE] 💾 Saved vibe [%d]: %s\n", newVibe.ID, newVibe.Message)
}

// handleList displays all saved vibes.
func handleList() {
	vibes := readVibes()
	if len(vibes) == 0 {
		fmt.Println("[VIBE] No vibes saved yet. Use 'vibe save \"<message>\"' to save one.")
		return
	}

	fmt.Println("[VIBE] --- VIBE HISTORY ---")
	for i := len(vibes) - 1; i >= 0; i-- {
		v := vibes[i]
		fmt.Printf("ID: %d | %s | %s\n", v.ID, v.Timestamp, v.Message)
	}
	fmt.Println("[VIBE] --------------------")
}

// handleCheckout reverts the working directory to a saved vibe.
func handleCheckout(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vibe checkout <vibe_id>")
		return
	}
	vibeID := args[0]
	patchFile := fmt.Sprintf("%s/%s.patch", patchesDir, vibeID)

	if _, err := os.Stat(patchFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "[VIBE] Error: Vibe with ID %s not found.\n", vibeID)
		os.Exit(1)
	}

	if err := exec.Command("git", "checkout", "--", ".").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error cleaning workspace: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("git", "apply", "--3way", patchFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error applying patch:\n%s\n", output)
		os.Exit(1)
	}

	fmt.Printf("[VIBE] ⏪ Workspace reverted to vibe [%s].\n", vibeID)
}

// handleDiff shows the contents of a specific patch file.
func handleDiff(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vibe diff <vibe_id>")
		return
	}
	vibeID := args[0]
	patchFile := fmt.Sprintf("%s/%s.patch", patchesDir, vibeID)

	data, err := os.ReadFile(patchFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error: Could not read vibe with ID %s.\n", vibeID)
		os.Exit(1)
	}

	fmt.Println(string(data))
}

// handleReset cleans the workspace and removes the .vibes directory.
func handleReset() {
	fmt.Println("[VIBE] 🧹 Cleaning workspace and removing session...")

	if err := exec.Command("git", "checkout", "--", ".").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error cleaning workspace: %v\n", err)
		os.Exit(1)
	}

	if err := os.RemoveAll(vibesDir); err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error removing .vibes directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("[VIBE] 💥 Vibegit session cleared.")
}

// --- Helper Functions ---

func readVibes() []Vibe {
	file, err := os.ReadFile(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "[VIBE] Error: vibegit not initialized. Run 'vibe init' first.")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "[VIBE] Error reading log file: %v\n", err)
		os.Exit(1)
	}

	var vibes []Vibe
	if err := json.Unmarshal(file, &vibes); err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error parsing log file: %v\n", err)
		os.Exit(1)
	}
	return vibes
}

func writeVibes(vibes []Vibe) {
	data, err := json.MarshalIndent(vibes, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error marshalling vibes: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(logFile, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error writing to log file: %v\n", err)
		os.Exit(1)
	}
}

