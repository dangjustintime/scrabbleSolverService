package main

import (
        "fmt"
        "github.com/gorilla/mux"
        "log"
        "net/http"
        "time"
)

const PORT = ":8080"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "Home Route")
}

func main() {
        router := mux.NewRouter()
        router.HandleFunc("/", HomeHandler)
        server := &http.Server{
                Handler: router,
                Addr: PORT,
                WriteTimeout: 10 * time.Second,
                ReadTimeout: 10 * time.Second,
        }

        log.Fatal(server.ListenAndServe())
}

