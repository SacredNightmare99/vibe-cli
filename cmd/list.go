package cmd

import "fmt"

func HandleList() {
	vibes := ReadVibes()
	if len(vibes) == 0 {
		fmt.Println("[VIBE] No vibes saved yet. Use `vibe save` to add one.")
		return
	}

	fmt.Println("[VIBE] --- VIBE HISTORY ---")
	for i := len(vibes) - 1; i >= 0; i-- {
		v := vibes[i]
		fmt.Printf("ID: %d | %s | %s\n", v.ID, v.Timestamp, v.Message)
	}
	fmt.Println("[VIBE] --------------------")
}

