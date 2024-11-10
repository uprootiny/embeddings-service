package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

// Updated OllamaResponse struct
type OllamaResponse struct {
    Model     string `json:"model"`
    CreatedAt string `json:"created_at"`
    Response  string `json:"response"`
    Done      bool   `json:"done"`
    // Include other fields if needed
}

// CallOllamaLLM sends a prompt to the Ollama LLM API and handles streaming response
func CallOllamaLLM(prompt string) (string, error) {
    url := "http://localhost:11434/api/generate"
    payload := map[string]interface{}{
        "model":  "llama3.2",
        "prompt": prompt,
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("error marshalling payload: %w", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("error creating request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    // Read the response line by line
    scanner := bufio.NewScanner(resp.Body)
    var fullResponse string

    for scanner.Scan() {
        line := scanner.Text()
        // Log the raw line for debugging
        log.Println("Ollama response line:", line)

        var ollamaResponse OllamaResponse
        if err := json.Unmarshal([]byte(line), &ollamaResponse); err != nil {
            // Handle error if necessary
            log.Println("Error unmarshalling line:", err)
            continue
        }

        fullResponse += ollamaResponse.Response
        if ollamaResponse.Done {
            break
        }
    }

    if err := scanner.Err(); err != nil {
        return "", fmt.Errorf("error reading response: %w", err)
    }

    return fullResponse, nil
}
