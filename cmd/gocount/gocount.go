package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "runtime"
    "strings"
)

var goRoutineLimit = 5
var total int
var sem = make(chan bool, goRoutineLimit)

func main() {

    filePath := "./data.txt"
    urls := getSliceOfUrls(filePath)

    fmt.Println("gr:", runtime.NumGoroutine())

    ch := make(chan int, len(urls))

    for _, v := range urls {
        sem <- true
        go fetch(v, ch)
    }

    for i := 0; i < len(urls); i++ {
        total += <-ch
    }
    fmt.Println("Total: ", total)

    fmt.Println("gr:", runtime.NumGoroutine())

}

func fetch(url string, ch chan<- int) {

    fmt.Println("gr:", runtime.NumGoroutine())
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error fetching url: ", url, err)
        return
    }
    defer resp.Body.Close()

    contents, err := ioutil.ReadAll(resp.Body)
    result := findGoMatches(contents)

    fmt.Println("count of ", url, result)
    ch <- result

    defer func() {
        <-sem
    }()

}

func getSliceOfUrls(path string) []string {
    result := []string{}

    file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        f := scanner.Text()
        result = append(result, f)
    }
    return result
}

func findGoMatches(s []byte) int {
    return strings.Count(string(s), "Go")
}
