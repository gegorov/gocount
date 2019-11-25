package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "runtime"
)

type urlCount struct {
    url       string
    countOfGo int
    body      string
}

func main() {

    filePath := "data.txt"
    urls := getSliceOfUrls(filePath)

    fmt.Println("gr:", runtime.NumGoroutine())

    ch := make(chan urlCount)

    go fetch(urls, ch)

    for v := range ch {
        fmt.Println(v.url, v.countOfGo)
    }
    fmt.Println("gr:", runtime.NumGoroutine())
    fmt.Println("Exiting")
}

func getSliceOfUrls(path string) []string {
    result := []string{}

    absPth, err := filepath.Abs(path)
    if err != nil {
        fmt.Println(err)
    }

    file, err := os.Open(absPth)
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

func fetch(urls []string, ch chan<- urlCount) {
    for _, url := range urls {
        fmt.Println("gr:", runtime.NumGoroutine())
        resp, err := http.Get(url)
        if err != nil {
            fmt.Println("Error fetching url: ", url, err)
            return
        }
        defer resp.Body.Close()

        contents, err := ioutil.ReadAll(resp.Body)
        result := urlCount{
            url:       url,
            countOfGo: findGoMatches(contents),
            body:      string(contents),
        }

        ch <- result
    }
    close(ch)

}

func findGoMatches(s []byte) int {
    re := regexp.MustCompile(`(?mi)go`)

    result := re.FindAllIndex(s, -1)

    return len(result)
}
