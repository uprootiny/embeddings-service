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
    "time"
)

// HomeHandler handles the root path and provides an operational dashboard
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    systemInfo := LoadSystemInfo()
    projects, err := GetRecentlyModifiedProjects()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error retrieving projects: %v", err), http.StatusInternalServerError)
        return
    }
    services := GetServiceStatus()
    scrapers := GetScraperStatus()

    htmlOutput := `
    <html>
    <head>
        <title>Operational Dashboard</title>
        <style>
            body { font-family: Arial, sans-serif; line-height: 1.6; margin: 20px; }
            .section { margin-bottom: 20px; }
            .highlight { background-color: #d4edda; padding: 5px; border-radius: 3px; }
            table { width: 100%; border-collapse: collapse; margin-top: 10px; }
            th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
            th { background-color: #f2f2f2; }
        </style>
    </head>
    <body>
        <h1>Operational Dashboard</h1>
        <div class="section">
            <h2>System Overview</h2>
            <p><strong>Hostname:</strong> %s</p>
            <p><strong>OS:</strong> %s</p>
            <p><strong>Uptime:</strong> %s</p>
            <p><strong>Kernel:</strong> %s</p>
            <p><strong>Architecture:</strong> %s</p>
        </div>
        <div class="section">
            <h2>Active Services</h2>
            <table>
                <thead>
                    <tr><th>Service</th><th>Port</th><th>Status</th><th>Actions</th></tr>
                </thead>
                <tbody>
    `

    // Populate services section
    for _, service := range services {
        statusClass := ""
        if service.Status == "Running" {
            statusClass = "highlight"
        }
        htmlOutput += fmt.Sprintf(
            `<tr>
                <td>%s</td>
                <td>%d</td>
                <td class="%s">%s</td>
                <td><a href="/view-logs?service=%s">View Logs</a></td>
             </tr>`,
            service.Name, service.Port, statusClass, service.Status, service.Name,
        )
    }

    htmlOutput += `
                </tbody>
            </table>
        </div>
        <div class="section">
            <h2>Recently Worked On Projects</h2>
            <table>
                <thead>
                    <tr><th>Project</th><th>Last Modified</th><th>Actions</th></tr>
                </thead>
                <tbody>
    `

    // Populate recent projects section
    for _, project := range projects {
        htmlOutput += fmt.Sprintf(
            `<tr>
                <td>%s</td>
                <td>%s</td>
                <td><a href="/repo-details?project=%s">Details</a></td>
             </tr>`,
            project.Name, project.LastModified, project.Name,
        )
    }

    htmlOutput += `
                </tbody>
            </table>
        </div>
        <div class="section">
            <h2>Scraper Management</h2>
            <table>
                <thead>
                    <tr><th>Scraper</th><th>Status</th><th>Actions</th></tr>
                </thead>
                <tbody>
    `

    // Populate scraper section
    for _, scraper := range scrapers {
        action := "Edit Config"
        htmlOutput += fmt.Sprintf(
            `<tr>
                <td>%s</td>
                <td>%s</td>
                <td>
                    <a href="/edit-config?scraper=%s">%s</a> |
                    <a href="/view-logs?scraper=%s">View Logs</a>
                </td>
             </tr>`,
            scraper.Name, scraper.Status, scraper.Name, action, scraper.Name,
        )
    }

    htmlOutput += `
                </tbody>
            </table>
        </div>
        <div class="section">
            <h2>Embeddings and Intent Mapping</h2>
            <p>Explore intents and the projects/tasks that fulfill them:</p>
            <ul>
                <li><strong>Intent:</strong> Scrape Financial News<br>
                    <strong>Matched Project:</strong> news_scraper<br>
                    <strong>Params:</strong> news_params.json</li>
                <li><strong>Intent:</strong> Perform Sentiment Analysis<br>
                    <strong>Matched Project:</strong> sentiment_analyzer<br>
                    <strong>Params:</strong> sentiment_params.json</li>
            </ul>
        </div>
    </body>
    </html>
    `

    // Fill in system information
    fmt.Fprintf(w, htmlOutput,
        systemInfo.Hostname, systemInfo.OS, systemInfo.Uptime,
        systemInfo.Kernel, systemInfo.Architecture,
    )
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

// GetRecentlyModifiedProjects scans directories and lists recently modified projects
func GetRecentlyModifiedProjects() ([]Repo, error) {
    basePaths := []string{"/home/uprootiny/ClojureProjects", "/home/uprootiny/Projects", "/home/uprootiny/tinystatus"}
    var repos []Repo

    for _, basePath := range basePaths {
        err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                log.Printf("Error accessing path %s: %v", path, err)
                return nil
            }

            if info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
                repos = append(repos, Repo{
                    Name:         info.Name(),
                    Path:         path,
                    LastModified: info.ModTime().Format("2006-01-02 15:04:05"),
                })
            }
            return nil
        })
        if err != nil {
            return nil, err
        }
    }

    // Sort repos by last modified date in descending order
    sort.Slice(repos, func(i, j int) bool {
        timeI, _ := time.Parse("2006-01-02 15:04:05", repos[i].LastModified)
        timeJ, _ := time.Parse("2006-01-02 15:04:05", repos[j].LastModified)
        return timeI.After(timeJ)
    })

    return repos, nil
}

// GetServiceStatus checks the status of active services
func GetServiceStatus() []ServiceStatus {
    services := []ServiceStatus{
        {"Ollama Server", 11435, "Not Running", "ollama-log.txt"},
        {"Electric App", 8085, "Running", "electric-log.txt"},
        {"Tinystatus", 4090, "Running", "tinystatus-log.txt"},
    }

    // Implement actual service status checks if necessary
    return services
}

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
                entryPoints := findEntryPoints(path)
                intents := inferIntents(path)

                repos = append(repos, Repo{
                    Name:         info.Name(),
                    Path:         path,
                    EntryPoints:  entryPoints,
                    Intents:      intents,
                    LastModified: info.ModTime().Format("2006-01-02 15:04:05"),
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
