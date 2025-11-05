package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	VibesDir    = ".vibes"
	PatchesDir  = ".vibes/patches"
	LogFile     = ".vibes/log.json"
	TrackedFile = ".vibes/tracked.json"
)

type Vibe struct {
	ID        int64  `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type Tracked struct {
	Projects []string `json:"projects"`
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

