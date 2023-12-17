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

    // Membuat klien HTTP dengan transport yang sudah disetel dengan proxy
    client := &http.Client{}

    // Membuat request dengan header khusus
    req, err := http.NewRequest("GET", targetURL, nil)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }

    req.Header.Add("User-Agent", "My-Go-App-User-Agent")
    req.Header.Add("X-Custom-Header", "CustomValue")

    rand.Seed(time.Now().UnixNano()) // Untuk mengacak proxy

    // Loop melalui daftar proxyURLs dan membuat request dengan menggunakan masing-masing proxy
    for _, proxyURL := range proxyURLs {
        // Mengacak urutan proxyURLs
        rand.Shuffle(len(proxyURLs), func(i, j int) {
            proxyURLs[i], proxyURLs[j] = proxyURLs[j], proxyURLs[i]
        })

        // Membuat transport dengan proxy yang digunakan dalam iterasi saat ini
        transport := &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        }

        // Mengatur transport dalam klien HTTP
        client.Transport = transport

        // Menjalankan request dengan klien HTTP yang sudah disiapkan
        resp, err := client.Do(req)
        if err != nil {
            fmt.Println("Error making request with proxy", proxyURL.String(), ":", err)
            continue
        }
        defer resp.Body.Close()

        // Mengambil body response atau melakukan operasi lainnya...
        fmt.Println("Request successful with proxy", proxyURL.String())

        // Menunda sebelum request berikutnya (misalnya, 5 detik)
        time.Sleep(500 * time.Nanosecond)
    }
}
