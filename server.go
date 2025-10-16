// server.go
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func callGeminiAPI(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	// Construct the JSON payload for the API request
	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %w", err)
	}

	// Make the HTTP POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Structs to parse the JSON response from Gemini
	type Part struct {
		Text string `json:"text"`
	}
	type Candidate struct {
		Content struct {
			Parts []Part `json:"parts"`
		} `json:"content"`
	}
	type APIResponse struct {
		Candidates []Candidate `json:"candidates"`
	}

	var apiResp APIResponse
	if err := json.Unmarshal(responseBody, &apiResp); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}
	
	// Extract the generated text from the response
	if len(apiResp.Candidates) > 0 && len(apiResp.Candidates[0].Content.Parts) > 0 {
		return apiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no content found in API response")
}

func runAndStream(ws *websocket.Conn, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		ws.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
	}
	return cmd.Wait()
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	log.Println("Client connected...")

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", err)
			break
		}
		rawPrompt := string(p)
		log.Printf("Received raw prompt: %s\n", rawPrompt)

		// 1. Parse the prompt into filename and instruction
		parts := strings.SplitN(rawPrompt, ":", 2)
		if len(parts) != 2 {
			ws.WriteMessage(websocket.TextMessage, []byte("[ERROR] Invalid prompt format. Use 'filename: instruction'"))
			continue
		}
		fileName := strings.TrimSpace(parts[0])
		userInstruction := strings.TrimSpace(parts[1])

		// 2. Read the content of the specified file
		ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[INFO] Reading file: %s...", fileName)))
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[ERROR] Could not read file '%s': %v", fileName, err)))
			continue
		}

		// 3. Construct a detailed prompt for the AI
		fullPrompt := fmt.Sprintf(`You are an expert AI programmer. Your task is to modify a code file based on a user's request.

File Name: %s

--- FILE CONTENT START ---
%s
--- FILE CONTENT END ---

User's Instruction: %s

Please provide ONLY the complete, new version of the file content as your response. Do not add any explanation, comments, or markdown formatting unless it is part of the code itself.`, fileName, string(fileContent), userInstruction)

		// 4. Call the Gemini API
		ws.WriteMessage(websocket.TextMessage, []byte("[INFO] Sending context to Gemini API..."))
		aiResponse, err := callGeminiAPI(fullPrompt)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("[ERROR] Gemini API failed: "+err.Error()))
			continue
		}
		
		// The AI might sometimes wrap the code in markdown, clean it.
		cleanedResponse := strings.TrimSpace(aiResponse)
		cleanedResponse = strings.TrimPrefix(cleanedResponse, "```go")
		cleanedResponse = strings.TrimPrefix(cleanedResponse, "```")
		cleanedResponse = strings.TrimSuffix(cleanedResponse, "```")
		cleanedResponse = strings.TrimSpace(cleanedResponse)


		// 5. Overwrite the original file with the AI's response
		ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[INFO] AI response received. Overwriting file %s...", fileName)))
		err = os.WriteFile(fileName, []byte(cleanedResponse), 0644) // 0644 is standard file permission
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[ERROR] Failed to write to file '%s': %v", fileName, err)))
			continue
		}

		// 6. Run vibe save to record the change
		ws.WriteMessage(websocket.TextMessage, []byte("[INFO] Saving vibe..."))
		// We use the original raw prompt for the vibe message
		err = runAndStream(ws, "./vibe", "save", rawPrompt)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("[ERROR] vibe save failed: "+err.Error()))
			continue
		}
		ws.WriteMessage(websocket.TextMessage, []byte("[SUCCESS] Workflow complete."))
	}
}

func startServer() {
	http.HandleFunc("/ws", handleConnection)
	fmt.Println("Vibe server listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


