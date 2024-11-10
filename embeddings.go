package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "strings"
)

// Embedding structure for intents and projects
type Embedding struct {
    Intent   string `json:"intent"`
    Project  string `json:"project"`
    Params   string `json:"params"`
}

// Load embeddings from a file
func LoadEmbeddings() ([]Embedding, error) {
    data, err := ioutil.ReadFile("data/embeddings.json")
    if err != nil {
        return nil, err
    }

    var embeddings []Embedding
    if err := json.Unmarshal(data, &embeddings); err != nil {
        return nil, err
    }
    return embeddings, nil
}

// MapIntentToProject maps a given intent to the most relevant project and parameters
func MapIntentToProject(intent string) Embedding {
    embeddings, err := LoadEmbeddings()
    if err != nil {
        log.Fatalf("Error loading embeddings: %v", err)
    }

    var bestMatch Embedding
    for _, embedding := range embeddings {
        if strings.Contains(strings.ToLower(intent), strings.ToLower(embedding.Intent)) {
            bestMatch = embedding
            break
        }
    }

    if bestMatch.Intent == "" {
        bestMatch = Embedding{
            Intent:  intent,
            Project: "No matching project found",
            Params:  "",
        }
    }
    return bestMatch
}
