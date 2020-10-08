package main

import (
        "encoding/json"
        "fmt"
        "github.com/gorilla/mux"
        "io"
        "log"
        "net/http"
        "os"
        "time"
)

const PORT = ":8080"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "Home Route")
}

func WordListHandler(w http.ResponseWriter, r *http.Request) {
        response, err := http.Get("http://recruiting.bluenile.com/words.txt")
        if err != nil {
                log.Fatal(err)
        }
        fmt.Fprint(w, response)
}

func WordsHandler(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        params := mux.Vars(r)
        word, _ := params["word"]
        json.NewEncoder(w).Encode(word)
}

func DownloadWordList() {
        fileUrl := "http://recruiting.bluenile.com/words.txt"
        err := DownloadFile("words.txt", fileUrl)
        if err != nil {
                panic(err)
        }
        fmt.Println("Downloaded: " + fileUrl)
}

func DownloadFile(filepath string, url string) error {
        response, err := http.Get(url)
        if err != nil {
                return err
        }
        defer response.Body.Close()

        out, err := os.Create(filepath)
        if err != nil {
                return err
        }
        defer out.Close()
        _, err = io.Copy(out, response.Body)
        return err
}

func main() {
        router := mux.NewRouter()

        router.HandleFunc("/", HomeHandler)
        router.HandleFunc("/words/{word}", WordsHandler)
        router.HandleFunc("/wordlist", WordListHandler)

        server := &http.Server{
                Handler: router,
                Addr: PORT,
                WriteTimeout: 10 * time.Second,
                ReadTimeout: 10 * time.Second,
        }

        fmt.Print("server running...")
        DownloadWordList()
        log.Fatal(server.ListenAndServe())
}

