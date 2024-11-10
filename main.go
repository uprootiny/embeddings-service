package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", HomeHandler).Methods("GET")
    router.HandleFunc("/list-repos", ListReposHandler).Methods("GET")
    router.HandleFunc("/map-intent", MapIntentHandler).Methods("GET")
    router.HandleFunc("/repo-details", RepoDetailsHandler).Methods("GET")

    log.Println("Server running on port 8085")
    log.Fatal(http.ListenAndServe(":8085", router))
}

// HomeHandler handles the root path
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to the Embeddings Service API"))
}
