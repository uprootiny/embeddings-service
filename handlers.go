package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os/exec"
    "strings"
    "log"
)

// Repo represents a basic project structure with relevant data
type Repo struct {
    Name         string `json:"name"`
    LastModified string `json:"last_modified"`
    Path         string `json:"path"`
}

// ListRecentRepos lists recently modified project directories
func ListRecentRepos() []Repo {
    var repos []Repo
    projectPaths := []string{"/home/uprootiny/ClojureProjects", "/home/uprootiny/Projects/November"}

    findCmd := fmt.Sprintf("find %s -maxdepth 1 -type d -printf '%%TY-%%Tm-%%Td %%TT %%p\\n' | sort -r | head -n 5", strings.Join(projectPaths, " "))
    out, err := exec.Command("bash", "-c", findCmd).Output()
    if err != nil {
        log.Printf("Error listing repos: %v", err)
        return repos
    }

    for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
        if line == "" {
            continue
        }
        parts := strings.SplitN(line, " ", 3)
        if len(parts) == 3 {
            repos = append(repos, Repo{
                LastModified: parts[0] + " " + parts[1],
                Path:         parts[2],
                Name:         parts[2], // Customize this as needed
            })
        }
    }
    return repos
}

// ListReposHandler is an HTTP handler for listing recent repos
func ListReposHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling request to list recent repos")
    repos := ListRecentRepos()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(repos)
}

// MapIntentHandler maps user intents to project invocations
func MapIntentHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling request to map user intent")
    // Placeholder logic for mapping intents; this should use embeddings.
    intents := []string{"Scrape Financial News", "Run Sentiment Analysis"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(intents)
}

// RepoDetailsHandler provides details for a specific repo
func RepoDetailsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling request for repo details")
    repoPath := r.URL.Query().Get("path")
    if repoPath == "" {
        http.Error(w, "Missing 'path' parameter", http.StatusBadRequest)
        return
    }

    // Example command to list files in a directory
    cmd := exec.Command("ls", "-l", repoPath)
    out, err := cmd.Output()
    if err != nil {
        http.Error(w, "Failed to retrieve repo details", http.StatusInternalServerError)
        log.Printf("Error retrieving repo details: %v", err)
        return
    }

    details := strings.Split(strings.TrimSpace(string(out)), "\n")
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(details)
}
