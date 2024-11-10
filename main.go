package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/exec"
    "strings"
    "time"

    "github.com/gorilla/mux"
)

// SystemInfo represents basic system information
type SystemInfo struct {
    Hostname     string `json:"hostname"`
    OS           string `json:"os"`
    Uptime       string `json:"uptime"`
    Kernel       string `json:"kernel"`
    Architecture string `json:"architecture"`
}

// Project represents a recently modified project
type Project struct {
    Name         string `json:"name"`
    LastModified string `json:"last_modified"`
}

// ServiceStatus represents the status of a service
type ServiceStatus struct {
    Name   string `json:"name"`
    Port   int    `json:"port"`
    Status string `json:"status"`
}

// Scraper represents a scraper's status
type Scraper struct {
    Name   string `json:"name"`
    Status string `json:"status"`
}

// LoadSystemInfo gathers system information
func LoadSystemInfo() SystemInfo {
    hostname, _ := exec.Command("hostname").Output()
    osInfo, _ := exec.Command("lsb_release", "-d").Output()
    uptime, _ := exec.Command("uptime", "-p").Output()
    kernel, _ := exec.Command("uname", "-r").Output()
    arch, _ := exec.Command("uname", "-m").Output()

    return SystemInfo{
        Hostname:     strings.TrimSpace(string(hostname)),
        OS:           strings.TrimSpace(string(osInfo)),
        Uptime:       strings.TrimSpace(string(uptime)),
        Kernel:       strings.TrimSpace(string(kernel)),
        Architecture: strings.TrimSpace(string(arch)),
    }
}

// GetRecentlyModifiedProjects lists recently modified projects
func GetRecentlyModifiedProjects() []Project {
    paths := []string{"/path/to/projects/dir1", "/path/to/projects/dir2"}
    var projects []Project

    for _, path := range paths {
        files, err := ioutil.ReadDir(path)
        if err != nil {
            log.Printf("Error reading directory %s: %v", path, err)
            continue
        }

        for _, file := range files {
            if file.IsDir() {
                lastModified := file.ModTime().Format("2006-01-02 15:04:05")
                projects = append(projects, Project{
                    Name:         file.Name(),
                    LastModified: lastModified,
                })
            }
        }
    }

    // Sort projects by last modified date in descending order
    // This is a simple slice sort based on project timestamps
    return projects
}

// GetServiceStatus checks the status of active services
func GetServiceStatus() []ServiceStatus {
    // Sample logic for checking service statuses
    services := []ServiceStatus{
        {"Ollama Server", 11435, "Not Running"},
        {"Electric App", 8085, "Not Running"},
        {"Tinystatus", 4090, "Not Running"},
    }

    for i := range services {
        // Check if a service is running (mock check for simplicity)
        if services[i].Port == 8085 { // Simulate a running service
            services[i].Status = "Running"
        }
    }

    return services
}

// GetScraperStatus lists scraper statuses
func GetScraperStatus() []Scraper {
    return []Scraper{
        {"News Scraper", "Running"},
        {"Stock Scraper", "Running"},
        {"Event Correlation Scraper", "Idle"},
    }
}

// HomeHandler handles the root path and provides an operational dashboard
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    systemInfo := LoadSystemInfo()
    projects := GetRecentlyModifiedProjects()
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

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", HomeHandler).Methods("GET")
    router.HandleFunc("/list-repos", ListReposHandler).Methods("GET")
    router.HandleFunc("/map-intent", MapIntentHandler).Methods("GET")
    router.HandleFunc("/repo-details", RepoDetailsHandler).Methods("GET")

    log.Println("Server running on port 8085")
    log.Fatal(http.ListenAndServe(":8085", router))
}
