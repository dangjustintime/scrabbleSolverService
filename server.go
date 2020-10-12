package main

import (
        "bufio"
        "encoding/json"
        "fmt"
        "github.com/gorilla/mux"
        "io"
        "log"
        "net/http"
        "os"
        "time"
)

// variables
const FILENAME = "words.txt"
const PORT = ":8080"
const TIMEOUTTIME = 10 * time.Second
var SCRABBLELETTERS = map[string]int{
        "a": 1, "e": 1, "i": 1, "l": 1, "n": 1, "o": 1, "r": 1, "s": 1, "t": 1, "u": 1,
        "d": 2, "g": 2,
        "b": 3, "c": 3, "m": 3, "p": 3,
        "f": 4, "h": 4, "v": 4, "w": 4, "y": 4,
        "k": 5,
        "j": 8, "x": 8,
        "q": 10, "z": 10,
}

// file functions
func FileExists(filename string) bool {
        info, err := os.Stat(filename)
        if os.IsNotExist(err) {
                return false
        }
        return !info.IsDir()
}

func DownloadWordList() {
        fileUrl := "http://recruiting.bluenile.com/words.txt"
        err := DownloadFile(FILENAME, fileUrl)
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

func ReadFile(fileName string) map[string]int {
        wordsMap := make(map[string]int)

        f, err := os.Open(fileName)
        if err != nil {
                log.Fatal(err)
        }
        defer f.Close()

        scanner := bufio.NewScanner(f)

        for scanner.Scan() {
                wordsMap[scanner.Text()] = 1
        }

        if err := scanner.Err(); err != nil {
                log.Fatal(err)
        }
        return wordsMap
}

// helper functions
func GetScrabbleScore(word string) int {
        var score int = 0
        for i := 0; i < len(word); i++ {
               score += SCRABBLELETTERS[string(word[i])]
        }
        return score
}

func GetWords(letters []string, wordsMap map[string]int) []string {
        var words []string = []string{}
        for i := 0; i < len(letters); i++ {
                var newLetters []string
                newLetters = append(newLetters, letters...)
                newLetters[i] = newLetters[len(newLetters) - 1]
                newLetters[len(newLetters) - 1] = ""
                newLetters = newLetters[:len(newLetters) - 1]
                var combinations []string
                GetCombinations(letters[i], &combinations, newLetters, wordsMap)
                for j := range combinations {
                        var hasWord bool = false
                        for k := range words {
                                if words[k] == combinations[j] {
                                        hasWord = true
                                }
                        }
                        if !hasWord {
                                words = append(words, combinations[j])
                        }
                }
        }
        return words
}

func MergeSort(words []string) []string {
        if len(words) <= 1 {
                return words
        }
        mid := len(words) / 2
        left := MergeSort(words[:mid])
        right := MergeSort(words[mid:])
        var wordsCopy []string
        wordsCopy = append(wordsCopy, words...)
        return Merge(left, right, wordsCopy)
}

func Merge(left []string, right []string, words []string) []string {
        leftCursor := 0
        rightCursor := 0
        for leftCursor < len(left) && rightCursor < len(right) {
                if GetScrabbleScore(left[leftCursor]) > GetScrabbleScore(right[rightCursor]) {
                        words[leftCursor + rightCursor] = left[leftCursor]
                        leftCursor++
                } else {
                        words[leftCursor + rightCursor] = right[rightCursor]
                        rightCursor++
                }
        }
        for i := leftCursor; i < len(left); i++ {
                words[leftCursor + rightCursor] = left[leftCursor]
        }
        for i := rightCursor; i < len(right); i++ {
                words[leftCursor + rightCursor] = right[rightCursor]
        }
        return words
}

func GetCombinations(word string, words *[]string, letters []string, wordsMap map[string]int) {
        if _, ok := wordsMap[word]; ok {
                (*words) = append((*words), word)
        }
        for i:= 0; i < len(letters); i++ {
                var newLetters []string
                newLetters = append(newLetters, letters...)
                newLetters[i] = newLetters[len(newLetters) - 1]
                newLetters[len(newLetters) - 1] = ""
                newLetters = newLetters[:len(newLetters) - 1]
                newWord := word + letters[i]
                GetCombinations(newWord, words, newLetters, wordsMap)
        }
}

// handlers
func HomeHandler(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "Home Route")
}

func WordsHandler(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        params := mux.Vars(r)
        word, _ := params["word"]
        var letters []string
        for i := 0; i < len(word); i++ {
                letters = append(letters, string(word[i]))
        }

        if !FileExists(FILENAME) {
                DownloadWordList()
        }
        words := GetWords(letters, ReadFile(FILENAME))
        json.NewEncoder(w).Encode(words)
}

// main
func main() {
        router := mux.NewRouter()

        router.HandleFunc("/", HomeHandler)
        router.HandleFunc("/words/{word}", WordsHandler)

        server := &http.Server{
                Handler: router,
                Addr: PORT,
                WriteTimeout: TIMEOUTTIME,
                ReadTimeout: TIMEOUTTIME,
        }

        fmt.Println("\n\n\nserver running...")
        log.Fatal(server.ListenAndServe())
}
