package main

import (
        "encoding/json"
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

func WordsHandler(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        params := mux.Vars(r)
        word, _ := params["word"]
        json.NewEncoder(w).Encode(word)
}

func main() {
        router := mux.NewRouter()

        router.HandleFunc("/", HomeHandler)
        router.HandleFunc("/words/{word}", WordsHandler)

        server := &http.Server{
                Handler: router,
                Addr: PORT,
                WriteTimeout: 10 * time.Second,
                ReadTimeout: 10 * time.Second,
        }

        log.Fatal(server.ListenAndServe())
}

