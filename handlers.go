package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "sort"
    "strings"
    "html/template"

    "time"
)
func LLManalysisHandler(w http.ResponseWriter, r *http.Request) {
    ollamaPrompt := "Generate a brief analysis of the current state of services and active scrapers."
    ollamaResponse, err := CallOllamaLLM(ollamaPrompt)
    if err != nil {
        log.Println("Error calling Ollama LLM:", err)
        http.Error(w, "Error generating LLM response", http.StatusInternalServerError)
        return
    }

    response := map[string]string{"analysis": ollamaResponse}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// // HomeHandler handles the root path and provides an operational dashboard
// func HomeHandler(w http.ResponseWriter, r *http.Request) {
//     tmplPath := filepath.Join("templates", "dashboard.html")
//     tmpl, err := template.ParseFiles(tmplPath)
//     if err != nil {
//         http.Error(w, "Error loading template", http.StatusInternalServerError)
//         log.Println("Error parsing template:", err)
//         return
//     }

//     systemInfo := LoadSystemInfo()
//     projects, err := GetRecentlyModifiedProjects()
//     if err != nil {
//         http.Error(w, "Error retrieving projects", http.StatusInternalServerError)
//         log.Println("Error retrieving projects:", err)
//         return
//     }
//     services := GetServiceStatus()
//     scrapers := GetScraperStatus()

//     ollamaPrompt := "Generate a brief analysis of the current state of services and active scrapers."
//     ollamaResponse, err := CallOllamaLLM(ollamaPrompt)
//     if err != nil {
//         log.Println("Error calling Ollama LLM:", err)
//         ollamaResponse = "Error generating response from Ollama LLM."
//     }

//     intents := []IntentMapping{
//         {
//             Intent:         "Scrape Financial News",
//             MatchedProject: "news_scraper",
//             Params:         "news_params.json",
//         },
//         {
//             Intent:         "Perform Sentiment Analysis",
//             MatchedProject: "sentiment_analyzer",
//             Params:         "sentiment_params.json",
//         },
//     }

//     data := struct {
//         SystemInfo     SystemInfo
//         Services       []ServiceStatus
//         Projects       []Repo
//         Scrapers       []Scraper
//         Intents        []IntentMapping
//         OllamaResponse string
//     }{
//         SystemInfo:     systemInfo,
//         Services:       services,
//         Projects:       projects,
//         Scrapers:       scrapers,
//         Intents:        intents,
//         OllamaResponse: ollamaResponse,
//     }

//     err = tmpl.Execute(w, data)
//     if err != nil {
//         http.Error(w, "Error rendering template", http.StatusInternalServerError)
//         log.Println("Error executing template:", err)
//         return
//     }
// }
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    tmplPath := filepath.Join("templates", "dashboard.html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        log.Println("Error parsing template:", err)
        return
    }

    systemInfo := LoadSystemInfo()
    projects, err := GetRecentlyModifiedProjects()
    if err != nil {
        http.Error(w, "Error retrieving projects", http.StatusInternalServerError)
        log.Println("Error retrieving projects:", err)
        return
    }
    services := GetServiceStatus()
    scrapers := GetScraperStatus()

    ollamaPrompt := "Generate a brief analysis of the current state of services and active scrapers."
    ollamaResponse, err := CallOllamaLLM(ollamaPrompt)
    if err != nil {
        log.Println("Error calling Ollama LLM:", err)
        ollamaResponse = "Error generating response from Ollama LLM."
    }

    data := struct {
        SystemInfo     SystemInfo
        Services       []ServiceStatus
        Projects       []Repo
        Scrapers       []Scraper
        OllamaResponse string
    }{
        SystemInfo:     systemInfo,
        Services:       services,
        Projects:       projects,
        Scrapers:       scrapers,
        OllamaResponse: ollamaResponse,
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        log.Println("Error executing template:", err)
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}

// LoadSystemInfo gathers basic system information
func LoadSystemInfo() SystemInfo {
    hostname, _ := exec.Command("hostname").Output()
    osInfo, _ := exec.Command("lsb_release", "-d").Output()
    uptime, _ := exec.Command("uptime", "-p").Output()
    kernel, _ := exec.Command("uname", "-r").Output()
    arch, _ := exec.Command("uname", "-m").Output()

    return SystemInfo{
        Hostname:     strings.TrimSpace(string(hostname)),
        OS:           strings.TrimSpace(strings.Replace(string(osInfo), "Description:\t", "", 1)),
        Uptime:       strings.TrimSpace(string(uptime)),
        Kernel:       strings.TrimSpace(string(kernel)),
        Architecture: strings.TrimSpace(string(arch)),
    }
}

// func GetRecentlyModifiedProjects() ([]Repo, error) {
//     var repos []Repo

//     // Use environment variable or default paths
//     basePaths := []string{}

//     // Check if an environment variable is set for project paths
//     if pathsEnv := os.Getenv("PROJECT_PATHS"); pathsEnv != "" {
//         basePaths = strings.Split(pathsEnv, ":")
//     } else {
//         // Default paths
//         homeDir, err := os.UserHomeDir()
//         if err != nil {
//             return nil, fmt.Errorf("unable to get user home directory: %w", err)
//         }
//         basePaths = []string{
//             filepath.Join(homeDir, "Projects"),
//             filepath.Join(homeDir, "ClojureProjects"),
//             filepath.Join(homeDir, "tinystatus"),
//         }
//     }

//     seen := make(map[string]bool)

//     for _, basePath := range basePaths {
//         entries, err := ioutil.ReadDir(basePath)
//         if err != nil {
//             log.Printf("Error reading directory %s: %v", basePath, err)
//             continue
//         }

//         for _, entry := range entries {
//             // Skip if not a directory, already seen, or starts with a dot
//             if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") || seen[entry.Name()] {
//                 continue
//             }

//             repoPath := filepath.Join(basePath, entry.Name())

//             // Ensure we have permission to access the directory
//             if _, err := os.Stat(repoPath); err != nil {
//                 log.Printf("Cannot access directory %s: %v", repoPath, err)
//                 continue
//             }

//             repos = append(repos, Repo{
//                 Name:         entry.Name(),
//                 Path:         repoPath,
//                 LastModified: entry.ModTime(),
//             })
//             seen[entry.Name()] = true
//         }
//     }

//     // Sort repos by last modified date in descending order
//     sort.Slice(repos, func(i, j int) bool {
//         return repos[i].LastModified.After(repos[j].LastModified)
//     })

//     // Limit the number of projects displayed (e.g., top 10)
//     if len(repos) > 10 {
//         repos = repos[:10]
//     }

//     return repos, nil
// }
func GetRecentlyModifiedProjects() ([]Repo, error) {
    var repos []Repo

    // Determine the base paths for scanning
    basePaths := []string{}
    if pathsEnv := os.Getenv("PROJECT_PATHS"); pathsEnv != "" {
        basePaths = strings.Split(pathsEnv, ":")
    } else {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return nil, fmt.Errorf("unable to get user home directory: %w", err)
        }
        basePaths = []string{
            filepath.Join(homeDir, "Projects"),
            filepath.Join(homeDir, "ClojureProjects"),
            filepath.Join(homeDir, "tinystatus"),
            filepath.Join(homeDir, "NovProjects"), // Additional paths if needed
        }
    }

    seen := make(map[string]bool)

    for _, basePath := range basePaths {
        entries, err := ioutil.ReadDir(basePath)
        if err != nil {
            log.Printf("Error reading directory %s: %v", basePath, err)
            continue
        }

        for _, entry := range entries {
            if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") || seen[entry.Name()] {
                continue
            }

            repoPath := filepath.Join(basePath, entry.Name())

            if _, err := os.Stat(repoPath); err != nil {
                log.Printf("Cannot access directory %s: %v", repoPath, err)
                continue
            }

            repos = append(repos, Repo{
                Name:         entry.Name(),
                Path:         repoPath,
                LastModified: entry.ModTime(), // Use time.Time directly
            })
            seen[entry.Name()] = true
        }
    }

    // Sort repos by last modified date in descending order
    sort.Slice(repos, func(i, j int) bool {
        return repos[i].LastModified.After(repos[j].LastModified)
    })

    if len(repos) > 10 {
        repos = repos[:10]
    }

    return repos, nil
}

func GetServiceStatus() []ServiceStatus {
    services := []ServiceStatus{
        {"Ollama Server", 11435, "Not Running", "ollama-log.txt"},
        {"Electric App", 8085, "Running", "electric-log.txt"},
        {"Tinystatus", 4090, "Running", "tinystatus-log.txt"},
        {"News Scraper", 5000, "Running", "news-scraper-log.txt"},
        {"Stock Scraper", 5001, "Idle", "stock-scraper-log.txt"},
        {"Event Correlation Scraper", 5002, "Not Running", "event-correlation-log.txt"},
        {"Market Data Aggregator", 6000, "Running", "market-data-log.txt"},
        {"Sentiment Analysis Dashboard", 7000, "Running", "sentiment-dashboard-log.txt"},
        {"WebSocket Server", 9000, "Not Running", "websocket-log.txt"},
        {"Analytics Engine", 9100, "Running", "analytics-engine-log.txt"},
    }

    // Implement actual checks or integration with service discovery mechanisms if needed

    return services
}

// // GetServiceStatus checks the status of active services
// func GetServiceStatus() []ServiceStatus {
//     services := []ServiceStatus{
//         {"Ollama Server", 11435, "Not Running", "ollama-log.txt"},
//         {"Electric App", 8085, "Running", "electric-log.txt"},
//         {"Tinystatus", 4090, "Running", "tinystatus-log.txt"},
//     }

//     // Implement actual service status checks if necessary
//     return services
// }

// GetScraperStatus retrieves the status of active scrapers
func GetScraperStatus() []Scraper {
    scrapers := []Scraper{
        {"News Scraper", "Running", "/scrapers/news/config", "news-scraper-log.txt"},
        {"Stock Scraper", "Idle", "/scrapers/stock/config", "stock-scraper-log.txt"},
        {"Event Correlation Scraper", "Not Running", "/scrapers/correlation/config", ""},
    }

    // Implement actual scraper status logic if necessary
    return scrapers
}

// ListReposHandler handles requests for listing recent repositories
func ListReposHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling request to list recent repositories")
    basePaths := []string{"/home/uprootiny/ClojureProjects", "/home/uprootiny/Projects", "/home/uprootiny/tinystatus"}
    repos, err := ListWorkingDirs(basePaths)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error listing repos: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(repos)
}

// MapIntentHandler maps user intents to the most relevant project using embeddings
func MapIntentHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling request to map user intent")
    intent := r.URL.Query().Get("intent")
    if intent == "" {
        http.Error(w, "Missing 'intent' parameter", http.StatusBadRequest)
        return
    }

    embeddingsData, err := LoadEmbeddings("data/embeddings.json")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error loading embeddings: %v", err), http.StatusInternalServerError)
        return
    }

    bestMatch, similarity := MapIntentToProject(intent, embeddingsData)
    result := map[string]interface{}{
        "Intent":        intent,
        "MatchedProject": bestMatch.Project,
        "Params":         bestMatch.Params,
        "Similarity":     similarity,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

// RepoDetailsHandler provides detailed information about a specific repository
func RepoDetailsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling request for repo details")
    repoName := r.URL.Query().Get("project")
    if repoName == "" {
        http.Error(w, "Missing 'project' query parameter", http.StatusBadRequest)
        return
    }

    repoDetails, err := GetRepoDetails("/home/uprootiny", repoName)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting repo details: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(repoDetails)
}

// ListWorkingDirs scans directories and lists recently modified local working directories and repos
func ListWorkingDirs(basePaths []string) ([]Repo, error) {
    var repos []Repo

    for _, basePath := range basePaths {
        err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                log.Printf("Error accessing path %s: %v", path, err)
                return nil
            }

            if info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
                // Placeholder functions for finding entry points and inferring intents.
                entryPoints := findEntryPoints(path) // Implement this function as needed.
                intents := inferIntents(path)        // Implement this function as needed.

                repos = append(repos, Repo{
                    Name:         info.Name(),
                    Path:         path,
                    EntryPoints:  entryPoints,
                    Intents:      intents,
                    LastModified: info.ModTime(), // Use time.Time for LastModified.
                })
            }
            return nil
        })
        if err != nil {
            log.Printf("Error walking the path %s: %v", basePath, err)
        }
    }

    return repos, nil
}


// GetRepoDetails retrieves detailed information about a specific repository
func GetRepoDetails(basePath, repoName string) (map[string]string, error) {
    repoPath := filepath.Join(basePath, repoName)
    if _, err := os.Stat(repoPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("repository not found")
    }

    details := map[string]string{
        "Name":     repoName,
        "Path":     repoPath,
        "Modified": time.Now().Format("2006-01-02 15:04:05"),
    }

    gitStatusCmd := exec.Command("git", "-C", repoPath, "status", "--short")
    output, err := gitStatusCmd.Output()
    if err != nil {
        return nil, fmt.Errorf("error running 'git status': %v", err)
    }
    details["Status"] = strings.TrimSpace(string(output))

    return details, nil
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
        if !file.IsDir() && (strings.HasSuffix(file.Name(), ".go") || strings.HasSuffix(file.Name(), ".clj") ||
            strings.HasSuffix(file.Name(), ".js") || strings.HasSuffix(file.Name(), ".py") || strings.HasSuffix(file.Name(), ".sh")) {
            entryPoints = append(entryPoints, file.Name())
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
