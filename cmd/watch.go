package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func HandleWatch() {
	data, err := os.ReadFile(TrackedFile)
	if err != nil {
		fmt.Println("[VIBE] No tracked projects found. Run `vibe init` first.")
		return
	}

	var tracked Tracked
	if err := json.Unmarshal(data, &tracked); err != nil {
		fmt.Println("[VIBE] Error parsing tracked projects:", err)
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, dir := range tracked.Projects {
		fmt.Println("[VIBE] 👁️ Watching:", dir)
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err == nil && info.IsDir() {
				_ = watcher.Add(path)
			}
			return nil
		})
	}

	lastChange := time.Now()
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				if time.Since(lastChange) > 4*time.Second {
					fmt.Println("[VIBE] 🔍 Change detected:", event.Name)
					lastChange = time.Now()
					saveChange()
				}
			}
		case err := <-watcher.Errors:
			fmt.Println("[VIBE] Watcher error:", err)
		}
	}
}

func saveChange() {
	fmt.Println("[VIBE] 💾 Auto-saving vibe...")
	cmd := exec.Command("./vibe", "save", "Auto-save: Gemini edit")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	fmt.Println("[VIBE] ✅ Auto-save complete.")
}

