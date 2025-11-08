package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Paths
const (
	VibesDir    = ".vibes"
	PatchesDir  = ".vibes/patches"
	LogFile     = ".vibes/log.json"
	CentralTrackedFile = "tracked.json"
	CentralVibeDir = ".vibe"
)

type Vibe struct {
	ID        int64  `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type Project struct {
	ID string `json:"id"`
	Path string `json:"path"`
}

type Tracked struct {
	Projects []Project `json:"projects"`
}

// Returns the absolute path to the central tracked.json
func GetCentralTrackedFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, CentralVibeDir)
	return filepath.Join(configDir, CentralTrackedFile), nil
}

func ReadTrackedProjects() (Tracked, error) {
	var tracked Tracked
	trackedFile, err := GetCentralTrackedFile()
	if err != nil {
		return tracked, fmt.Errorf("could not get home dir: %w", err)
	}

	data, err := os.ReadFile(trackedFile)
	if err != nil {
		if os.IsNotExist(err) {
			return Tracked{}, nil
		}
		return tracked, err
	}

	if err := json.Unmarshal(data, &tracked); err != nil {
		return tracked, err
	}
	return tracked, nil
}

func WriteTrackedProjects(tracked Tracked) error {
	trackedFile, err := GetCentralTrackedFile()
	if err != nil {
		return fmt.Errorf("could not get home dir: %w", err)
	}

	// Ensure the central .vibe directory exists
	configDir := filepath.Dir(trackedFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create central config dir: %w", err)
	}

	data, _ := json.MarshalIndent(tracked, "", "  ")
	return os.WriteFile(trackedFile, data, 0644)
}

// ReadVibes loads the JSON log file.
func ReadVibes() []Vibe {
	data, err := os.ReadFile(LogFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("[VIBE] Not initialized. Run `vibe init` first.")
			os.Exit(1)
		}
		panic(err)
	}
	var vibes []Vibe
	_ = json.Unmarshal(data, &vibes)
	return vibes
}

// WriteVibes writes the vibe list back to log.json.
func WriteVibes(vibes []Vibe) {
	data, _ := json.MarshalIndent(vibes, "", "  ")
	_ = os.WriteFile(LogFile, data, 0644)
}

