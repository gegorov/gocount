package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
    "sync"
)

const goRoutineLimit = 5

func main() {
    var total int

    ch := make(chan int)

    filePath := "./data.txt"
    urls, err := getSliceOfUrls(filePath)
    if err != nil {
        log.Fatal(err)
    }

    go producer(urls, ch)

    for count := range ch {
        total += count
    }
    fmt.Println("Total: ", total)

}

func producer(urls []string, ch chan int) {
    var wg sync.WaitGroup
    sem := make(chan bool, goRoutineLimit)

    for _, v := range urls {
        url := v
        sem <- true
        wg.Add(1)
        go func() {
            fetch(url, ch)
            <-sem
            wg.Done()
        }()
    }

    wg.Wait()
    close(ch)
}

func fetch(url string, ch chan<- int) {

    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error fetching url: ", url, err)
        return
    }
    defer resp.Body.Close()

    contents, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("can't read body (%s): %s", url, err)
        return
    }

    result := strings.Count(string(contents), "Go")

    fmt.Println("count of ", url, result)
    ch <- result
}

func getSliceOfUrls(path string) ([]string, error) {
    result := []string{}

    file, err := os.Open(path)
    if err != nil {
        return result, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        url := scanner.Text()
        result = append(result, url)
    }
    return result, nil
}
