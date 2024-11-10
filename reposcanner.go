package main

import (
    "os"
    "path/filepath"
    "log"
    "strings"
)

// ScanDirectories scans the given paths for project directories
func ScanDirectories(basePaths []string) []Repo {
    var repos []Repo
    for _, basePath := range basePaths {
        err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                log.Printf("Error accessing path %s: %v", path, err)
                return nil
            }
            if info.IsDir() && path != basePath && !strings.HasPrefix(info.Name(), ".") {
                repos = append(repos, Repo{
                    Name: info.Name(),
                    Path: path,
                })
            }
            return nil
        })
        if err != nil {
            log.Printf("Error walking the path %s: %v", basePath, err)
        }
    }
    return repos
}