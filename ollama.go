package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

func main() {
    url := "http://localhost:11434/api/generate"
    payload := map[string]interface{}{
        "model":  "llama3.2",
        "prompt": "Why is the sky blue?",
    }
    jsonData, err := json.Marshal(payload)
    if err != nil {
        log.Fatalf("Error marshalling payload: %v", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        log.Fatalf("Error creating request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Error sending request: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Error reading response body: %v", err)
    }

    fmt.Println("Response:", string(body))
}
