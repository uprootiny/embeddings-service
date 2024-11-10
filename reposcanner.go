package main

import (
    "os"
    "path/filepath"
    "strings"
)

// Repo holds information about a project repository
type Repo struct {
    Name         string   `json:"name"`
    Path         string   `json:"path"`
    EntryPoints  []string `json:"entryPoints"`
    Intents      []string `json:"intents"`
}

// ListRecentRepos lists recently modified project directories
func ListRecentRepos() []Repo {
    var repos []Repo
    projectPaths := []string{"/home
