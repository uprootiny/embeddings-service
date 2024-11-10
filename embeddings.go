package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "math"
    "os"
    "strings"
    "unicode"
)

// Embedding structure for intents and projects
type Embedding struct {
    Intent   string    `json:"intent"`
    Project  string    `json:"project"`
    Params   string    `json:"params"`
    Vector   []float64 `json:"vector"` // Add a vector field for similarity matching
}

// Word embeddings map loaded from an external file
var wordEmbeddings = make(map[string][]float64)
var embeddingDimension = 4 // Change this based on your embedding file's dimension

// Load word embeddings from a file
func LoadWordEmbeddings(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("failed to open embeddings file: %v", err)
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&wordEmbeddings); err != nil {
        return fmt.Errorf("failed to decode embeddings: %v", err)
    }

    // Set the embedding dimension based on the first vector in the map
    for _, vec := range wordEmbeddings {
        embeddingDimension = len(vec)
        break
    }
    log.Printf("Loaded %d word embeddings with dimension %d", len(wordEmbeddings), embeddingDimension)
    return nil
}

// Load embeddings from a file
func LoadEmbeddings(filePath string) ([]Embedding, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read embeddings file: %v", err)
    }

    var embeddings []Embedding
    if err := json.Unmarshal(data, &embeddings); err != nil {
        return nil, fmt.Errorf("failed to unmarshal embeddings: %v", err)
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

// tokenize splits a sentence into words, removing punctuation and converting to lowercase
func tokenize(sentence string) []string {
    var tokens []string
    words := strings.FieldsFunc(sentence, func(r rune) bool {
        return !unicode.IsLetter(r) && !unicode.IsNumber(r)
    })

    for _, word := range words {
        tokens = append(tokens, strings.ToLower(word))
    }

    return tokens
}

// getWordEmbedding retrieves the embedding for a word, or a zero vector if the word is not found
func getWordEmbedding(word string) []float64 {
    if embedding, exists := wordEmbeddings[word]; exists {
        return embedding
    }
    // Return a zero vector if the word is not in the lookup table
    zeroVector := make([]float64, embeddingDimension)
    return zeroVector
}

// averageVectors averages a list of vectors to create a single sentence vector
func averageVectors(vectors [][]float64) []float64 {
    if len(vectors) == 0 {
        return make([]float64, embeddingDimension)
    }

    result := make([]float64, embeddingDimension)
    for _, vector := range vectors {
        for i, value := range vector {
            result[i] += value
        }
    }

    for i := range result {
        result[i] /= float64(len(vectors))
    }

    return result
}

// convertIntentToVector converts an intent into a sentence vector
func convertIntentToVector(intent string) []float64 {
    tokens := tokenize(intent)
    var vectors [][]float64

    for _, token := range tokens {
        vectors = append(vectors, getWordEmbedding(token))
    }

    sentenceVector := averageVectors(vectors)
    log.Printf("Generated vector for intent '%s': %v", intent, sentenceVector)

    return sentenceVector
}

// MapIntentToProject maps a given intent to the most relevant project using embeddings similarity
func MapIntentToProject(intent string, embeddings []Embedding) (Embedding, float64) {
    intentVector := convertIntentToVector(intent)
    var bestMatch Embedding
    highestSimilarity := -1.0

    for _, embedding := range embeddings {
        similarity := CalculateCosineSimilarity(intentVector, embedding.Vector)
        if similarity > highestSimilarity {
            highestSimilarity = similarity
            bestMatch = embedding
        }
    }

    return bestMatch, highestSimilarity
}

func main() {
    // Load word embeddings from an external file
    if err := LoadWordEmbeddings("data/word_embeddings.json"); err != nil {
        log.Fatalf("Error loading word embeddings: %v", err)
    }

    // Load intent embeddings
    embeddings, err := LoadEmbeddings("data/embeddings.json")
    if err != nil {
        log.Fatalf("Error loading intent embeddings: %v", err)
    }

    // Test mapping an intent to a project
    intent := "Analyze market trends and news"
    bestMatch, similarity := MapIntentToProject(intent, embeddings)

    log.Printf("Best match for intent '%s': Project: %s, Params: %s, Similarity: %f",
        intent, bestMatch.Project, bestMatch.Params, similarity)
}
