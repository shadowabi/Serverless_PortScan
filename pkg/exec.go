package pkg

import (
    "os"
    "bufio"
    "strings"
    "fmt"
    "sync"
    "regexp"
    "net/http"
    "io/ioutil"
    "time"
    "sort"
    "encoding/json"
    "crypto/tls"
)

var (
    rs []string // 存放结果
    Url []*string // 存放正确Url，需传递
    Config Configure //存放配置文件，需传递
    mu sync.Mutex
    client = &http.Client{ //http.client
        Timeout: 5 * time.Second,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }
)

type Configure struct {
    Serverless   string `json:"Serverless"`
    Default_port string `json:"Default_port"`
}

func ReadConfig() {
    data, _ := ioutil.ReadFile("./config/config.json")

    // 解码 JSON 数据
    err := json.Unmarshal(data, &Config)
    if err != nil {
        fmt.Println("config.json配置出错!")
        os.Exit(1)
    }
    if Config.Serverless == "http://" {
        fmt.Println("请配置config.json")
        os.Exit(1)
    }
}

func CheckIP(ip string, wg *sync.WaitGroup) {
    defer wg.Done()

    // 当输入Url时提取出域名
    re := regexp.MustCompile(`(http|https)\:\/\/`)
    if re.MatchString(ip) {
        ip = re.ReplaceAllString(ip, "")
        re = regexp.MustCompile(`(\/|\\).*`)
        if re.MatchString(ip) {
            ip = re.ReplaceAllString(ip, "")
        }
    }

    // 匹配IP/域名
    if regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`).MatchString(ip) || regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,11}$`).MatchString(ip) {
        mu.Lock()
        Url = append(Url, &ip)
        mu.Unlock()
    }
}

func ReadFile(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("无法打开此文件")
        os.Exit(1)
    }
    scan := bufio.NewScanner(file)
    var wg sync.WaitGroup
    for scan.Scan() {
        line := strings.TrimSpace(scan.Text())
        wg.Add(1)
        go CheckIP(line, &wg)
        wg.Wait()
    }
    wg.Wait()
    file.Close()
}


func UniqueStrings(Urls []*string) []*string { //去重
    seen := make(map[string]bool)
    result := []*string{}
    for _, Url := range Urls {
        s := *Url
        if _, ok := seen[s]; !ok {
            seen[s] = true
            mu.Lock()
            result = append(result, Url)
            mu.Unlock()
        }
    }
    return result
}

func Scan(port_list string) {
    Url = UniqueStrings(Url)
    var wg sync.WaitGroup

    for _, ip := range Url {
        wg.Add(1)
        go func(ip string,wg *sync.WaitGroup) {
            defer wg.Done()
        
            target := Config.Serverless + "?ip=" + ip + "&port=" + port_list
            resp, err := client.Get(target)
            if err != nil {
                return
            }
            defer resp.Body.Close()

            if resp.Body != nil {
                body, _ := ioutil.ReadAll(resp.Body)
                if string(body) == `""` {
                    return
                }

                ports := strings.Split(strings.Trim(string(body),`"`), ",")
                for _, port := range ports {
                    fmt.Println(fmt.Sprintf("[+] %s:%s TCP OPEN.", ip, port))
                    mu.Lock()
                    rs = append(rs, ip + ":" + port)
                    mu.Unlock()
             
                }
            } else {
                return
            }
        }(*ip, &wg)
    }
    wg.Wait()
    OutPut()
}

func OutPut() {
    if len(rs) == 0 {
        fmt.Println("[-] 目标端口均已关闭")
        return
    }
    file, err := os.OpenFile("./result.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    defer writer.Flush()

    sort.Strings(rs)
    if rs != nil {
        for _, i := range rs {
            fmt.Fprintln(writer, i)
        }
    }

    fmt.Println("已保存到result.txt文件")
}