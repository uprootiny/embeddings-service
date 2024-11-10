package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "sort"
    "strings"
    "time"
    "io/ioutil"
)

// Repo holds information about a project or working directory
type Repo struct {
    Name        string   `json:"name"`
    Path        string   `json:"path"`
    EntryPoints []string `json:"entryPoints"`
    Intents     []string `json:"intents"`
    LastModTime string   `json:"lastModified"`
}

// MapIntentHandler handles requests to map user intents to project invocations
func MapIntentHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling request to map user intent")
    // Placeholder logic for mapping intents; replace with embeddings integration
    intents := []map[string]string{
        {"intent": "Scrape Financial News", "project": "news_scraper", "params": "news_params.json"},
        {"intent": "Run Sentiment Analysis", "project": "sentiment_analyzer", "params": "sentiment_params.json"},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(intents)
}

// ListReposHandler handles the request for listing recent repositories
func ListReposHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling request to list recent repositories")
    basePaths := []string{"/home/uprootiny/ClojureProjects", "/home/uprootiny/Projects", "/home/uprootiny/tinystatus"}
    repos := ListWorkingDirs(basePaths)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(repos)
}

// RepoDetailsHandler provides detailed information about a specific repo
func RepoDetailsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling request for repo details")
    repoPath := r.URL.Query().Get("path")
    if repoPath == "" {
        http.Error(w, "Missing 'path' parameter", http.StatusBadRequest)
        return
    }

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

// ListWorkingDirs scans directories and lists recently modified local working directories and repos
func ListWorkingDirs(basePaths []string) []Repo {
    var repos []Repo

    for _, basePath := range basePaths {
        err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                log.Printf("Error accessing path %s: %v", path, err)
                return nil
            }

            if info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
                entryPoints := findEntryPoints(path)
                intents := inferIntents(path)

                repos = append(repos, Repo{
                    Name:        info.Name(),
                    Path:        path,
                    EntryPoints: entryPoints,
                    Intents:     intents,
                    LastModTime: info.ModTime().Format(time.RFC3339),
                })
            }
            return nil
        })
        if err != nil {
            log.Printf("Error walking the path %s: %v", basePath, err)
        }
    }

    sort.Slice(repos, func(i, j int) bool {
        timeI, _ := time.Parse(time.RFC3339, repos[i].LastModTime)
        timeJ, _ := time.Parse(time.RFC3339, repos[j].LastModTime)
        return timeI.After(timeJ)
    })

    return repos
}

// findEntryPoints searches for common entry point files in a directory
func findEntryPoints(path string) []string {
    var entryPoints []string
    files, err := ioutil.ReadDir(path)
    if err != nil {
        log.Printf("Error reading directory %s: %v", path, err)
        return entryPoints
    }

    for _, file := range files {
        if !file.IsDir() {
            if strings.HasSuffix(file.Name(), ".go") || strings.HasSuffix(file.Name(), ".clj") ||
                strings.HasSuffix(file.Name(), ".js") || strings.HasSuffix(file.Name(), ".py") ||
                strings.HasSuffix(file.Name(), ".sh") {
                entryPoints = append(entryPoints, file.Name())
            }
        }
    }

    return entryPoints
}

// inferIntents attempts to deduce project or working directory intents
func inferIntents(path string) []string {
    var intents []string
    knownIntentFiles := []string{"README.md", "docs", "Makefile", "build.sh", "run.sh", "requirements.txt", "pom.xml", "Dockerfile"}

    for _, filename := range knownIntentFiles {
        if _, err := os.Stat(filepath.Join(path, filename)); err == nil {
            intents = append(intents, "Contains "+filename)
        }
    }

    if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
        intents = append(intents, "Git Repository")
    }

    return intents
}
