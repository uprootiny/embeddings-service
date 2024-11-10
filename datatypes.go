package main

import (
	"time"
)
// Shared data structures for the project

type SystemInfo struct {
    Hostname     string `json:"hostname"`
    OS           string `json:"os"`
    Uptime       string `json:"uptime"`
    Kernel       string `json:"kernel"`
    Architecture string `json:"architecture"`
}

type ServiceStatus struct {
    Name    string `json:"name"`
    Port    int    `json:"port"`
    Status  string `json:"status"`
    URL     string `json:"url,omitempty"` // Include this if needed for links
}



type Scraper struct {
    Name      string `json:"name"`
    Status    string `json:"status"`
    ConfigURL string `json:"configURL"`
    LogFile   string `json:"logFile"`
}

type Repo struct {
    Name         string   `json:"name"`
    Path         string   `json:"path"`
    EntryPoints  []string `json:"entryPoints"`
    Intents      []string `json:"intents"`
    LastModified time.Time `json:"lastModified"`
}

// type Repo struct {
//     Name         string
//     Path         string
//     LastModified  // Ensure LastModified is of type time.Time
// }


type Embedding struct {
    Intent   string    `json:"intent"`
    Project  string    `json:"project"`
    Params   string    `json:"params"`
    Vector   []float64 `json:"vector"`
}

type EmbeddingResult struct {
    Intent        string `json:"intent"`
    MatchedProject string `json:"matchedProject"`
    Params        string `json:"params"`
}

type DashboardData struct {
    SystemInfo SystemInfo        `json:"systemInfo"`
    Services   []ServiceStatus   `json:"services"`
    Projects   []Repo            `json:"projects"`
    Scrapers   []Scraper         `json:"scrapers"`
    Embeddings []EmbeddingResult `json:"embeddings"`
}


type IntentMapping struct {
    Intent         string
    MatchedProject string
    Params         string
}
