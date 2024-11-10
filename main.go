package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/map-intent", MapIntentHandler).Methods("GET")
    r.HandleFunc("/list-repos", ListReposHandler).Methods("GET")
    r.HandleFunc("/repo-details", RepoDetailsHandler).Methods("GET")

    port := ":8088" // Use a non-standard port
    log.Printf("Embeddings service running on port %s", port)
    log.Fatal(http.ListenAndServe(port, r))
}
