package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "math"
    "strings"
)

// Embedding structure for intents and projects
type Embedding struct {
    Intent   string    `json:"intent"`
    Project  string    `json:"project"`
    Params   string    `json:"params"`
    Vector   []float64 `json:"vector"` // Add a vector field for similarity matching
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

// CalculateCosineSimilarity calculates the cosine similarity between two vectors
func CalculateCosineSimilarity(vec1, vec2 []float64) float64 {
    if len(vec1) != len(vec2) {
        log.Println("Vectors are not the same length")
        return 0.0
    }

    var dotProduct, normVec1, normVec2 float64
    for i := range vec1 {
        dotProduct += vec1[i] * vec2[i]
        normVec1 += vec1[i] * vec1[i]
        normVec2 += vec2[i] * vec2[i]
    }

    if normVec1 == 0.0 || normVec2 == 0.0 {
        return 0.0
    }
    return dotProduct / (math.Sqrt(normVec1) * math.Sqrt(normVec2))
}

// MapIntentToProject maps a given intent to the most relevant project using embeddings similarity
func MapIntentToProject(intent string, intentVector []float64) Embedding {
    embeddings, err := LoadEmbeddings()
    if err != nil {
        log.Fatalf("Error loading embeddings: %v", err)
    }

    var bestMatch Embedding
    highestSimilarity := -1.0

    for _, embedding := range embeddings {
        similarity := CalculateCosineSimilarity(intentVector, embedding.Vector)
        if similarity > highestSimilarity {
            highestSimilarity = similarity
            bestMatch = embedding
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
