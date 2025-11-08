package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func HandleWatch(args []string) {
	tracked, err := ReadTrackedProjects()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[VIBE] Error reading central tracked projects: %v\n", err)
		return
	}

	if len(tracked.Projects) == 0 {
		fmt.Println("[VIBE] No tracked projects found. Run `vibe init` first.")
		return
	}

	var projectsToWatch []Project
	var projectFilter string

	if len(args) > 0 {
		projectFilter = args[0]
		found := false
		for _, p := range tracked.Projects {
			if p.ID == projectFilter {
				projectsToWatch = append(projectsToWatch, p)
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("[VIBE] Error: Project ID '%s' not found in tracked projects.\n", projectFilter)
			return
		}
	} else {
		projectsToWatch = tracked.Projects
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, p := range projectsToWatch {
		fmt.Println("[VIBE] 👁️ Watching:", p.ID, "(", p.Path, ")")
		// Walk the directory and add all subdirectories to the watcher
		filepath.Walk(p.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				// Don't watch the .vibes or .git directories
				if info.Name() == VibesDir || info.Name() == ".git" {
					return filepath.SkipDir
				}
				if err := watcher.Add(path); err != nil {
					log.Println("[VIBE] Error adding path to watcher:", err)
				}
			}
			return nil
		})
	}

	// Use a map to debounce change events per project
	lastChanges := make(map[string]time.Time)

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				var project Project
				found := false
				for _, p := range projectsToWatch {
					// Find which project this event belongs to
					if strings.HasPrefix(event.Name, p.Path) {
						project = p
						found = true
						break
					}
				}

				if found {
					// Debounce: only save if last change for this project was > 4s ago
					if time.Since(lastChanges[project.ID]) > 4*time.Second {
						fmt.Printf("[VIBE] 🔍 Change detected in %s: %s\n", project.ID, event.Name)
						lastChanges[project.ID] = time.Now()
						// Run save in a goroutine so it doesn't block the watcher
						go saveChange(project.Path)
					}
				}
			}
		case err := <-watcher.Errors:
			fmt.Println("[VIBE] Watcher error:", err)
		}
	}
}

// saveChange runs the 'vibe save' command in the specified project directory
func saveChange(projectPath string) {
	// Get the path to the currently running 'vibe' executable
	vibePath, err := os.Executable()
	if err != nil {
		log.Printf("[VIBE] Error finding vibe executable for project %s: %v\n", projectPath, err)
		return
	}

	fmt.Printf("[VIBE] 💾 Auto-saving vibe for %s...\n", projectPath)
	cmd := exec.Command(vibePath, "save", "Auto-save: Gemini edit")
	cmd.Dir = projectPath 
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("[VIBE] Error auto-saving project %s: %v\n", projectPath, err)
	}
	fmt.Printf("[VIBE] ✅ Auto-save complete for %s.\n", projectPath)
}
