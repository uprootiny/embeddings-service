package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    // Create a new router
    router := mux.NewRouter()

    // Register route handlers
    router.HandleFunc("/", HomeHandler).Methods("GET")
    router.HandleFunc("/list-repos", ListReposHandler).Methods("GET")
    router.HandleFunc("/map-intent", MapIntentHandler).Methods("GET")
    router.HandleFunc("/repo-details", RepoDetailsHandler).Methods("GET")

    // Start the server
    log.Println("Server running on port 8085")
    log.Fatal(http.ListenAndServe(":8085", router))
}
