package main

import (
    "os"
    // "path/filepath"
    // "log"
    // "io/ioutil"
    // "strings"
    // "time"
)

// // findEntryPoints searches for common entry point files in a directory
// func findEntryPoints(path string) []string {
//     var entryPoints []string
//     files, err := ioutil.ReadDir(path)
//     if err != nil {
//         log.Printf("Error reading directory %s: %v", path, err)
//         return entryPoints
//     }

//     for _, file := range files {
//         if !file.IsDir() && (strings.HasSuffix(file.Name(), ".go") || strings.HasSuffix(file.Name(), ".clj") ||
//             strings.HasSuffix(file.Name(), ".js") || strings.HasSuffix(file.Name(), ".py") || strings.HasSuffix(file.Name(), ".sh")) {
//             entryPoints = append(entryPoints, file.Name())
//         }
//     }

//     return entryPoints
// }

// // inferIntents attempts to deduce project or working directory intents
// func inferIntents(path string) []string {
//     var intents []string
//     knownIntentFiles := []string{"README.md", "docs", "Makefile", "build.sh", "run.sh", "requirements.txt", "pom.xml", "Dockerfile"}

//     for _, filename := range knownIntentFiles {
//         if _, err := os.Stat(filepath.Join(path, filename)); err == nil {
//             intents = append(intents, "Contains "+filename)
//         }
//     }

//     if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
//         intents = append(intents, "Git Repository")
//     }

//     return intents
// }

// Utility function to format the last modified time of a file
func formatLastModifiedTime(info os.FileInfo) string {
    return info.ModTime().Format("2006-01-02 15:04:05")
}
