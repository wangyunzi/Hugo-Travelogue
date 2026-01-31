package main

import (
  "bufio"
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
  "net/url"
  "os"
  "path/filepath"
  "sort"
  "strings"
  "sync"
  "time"

  "github.com/mmcdole/gofeed"
)

const (
  feedsPath   = "rss/rss_feeds.txt"
  avatarsPath = "data/avatar_data.json"
  outputPath  = "data/rss_data.json"
  logPath     = "logs/error.log"

  maxRetries    = 3
  retryInterval = 10 * time.Second
)

type Avatar struct {
  Name   string `json:"name"`
  Avatar string `json:"avatar"`
}

type Article struct {
  DomainName string    `json:"domainName"`
  Name       string    `json:"name"`
  Title      string    `json:"title"`
  Link       string    `json:"link"`
  Date       string    `json:"date"`
  Avatar     string    `json:"avatar"`
  DateTime   time.Time `json:"-"`
}

func logError(format string, args ...interface{}) {
  message := fmt.Sprintf(format, args...)
  timestamp := time.Now().Format("2006-01-02 15:04:05")
  line := fmt.Sprintf("[%s] %s\n", timestamp, message)

  if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
    fmt.Printf("error creating log dir: %v\n", err)
    return
  }

  f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    fmt.Printf("error opening log file: %v\n", err)
    return
  }
  defer f.Close()

  if _, err := f.WriteString(line); err != nil {
    fmt.Printf("error writing log: %v\n", err)
  }
}

func cleanXMLContent(content string) string {
  var buf bytes.Buffer
  for _, r := range content {
    if r == '\t' || r == '\n' || r == '\r' || r >= 0x20 {
      buf.WriteRune(r)
    }
  }
  return buf.String()
}

func extractDomain(urlStr string) (string, error) {
  u, err := url.Parse(urlStr)
  if err != nil {
    return "", err
  }
  domain := u.Hostname()
  protocol := "https://"
  if u.Scheme != "" {
    protocol = u.Scheme + "://"
  }
  if domain == "" {
    return "", fmt.Errorf("empty domain from url: %s", urlStr)
  }
  return protocol + domain, nil
}

func readFeeds(path string) ([]string, error) {
  data, err := os.ReadFile(path)
  if err != nil {
    return nil, err
  }

  var feeds []string
  scanner := bufio.NewScanner(bytes.NewReader(data))
  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if line == "" || strings.HasPrefix(line, "#") {
      continue
    }
    feeds = append(feeds, line)
  }
  if err := scanner.Err(); err != nil {
    return nil, err
  }
  return feeds, nil
}

func loadAvatars(path string) (map[string]string, error) {
  data, err := os.ReadFile(path)
  if err != nil {
    return nil, err
  }

  var avatars []Avatar
  if err := json.Unmarshal(data, &avatars); err != nil {
    return nil, err
  }

  avatarMap := make(map[string]string)
  for _, a := range avatars {
    avatarMap[a.Name] = a.Avatar
  }
  return avatarMap, nil
}

func fetchRSS(feeds []string, avatars map[string]string) ([]Article, error) {
  var articles []Article
  var mu sync.Mutex
  var wg sync.WaitGroup

  fp := gofeed.NewParser()
  httpClient := &http.Client{Timeout: 10 * time.Second}

  for _, feedURL := range feeds {
    wg.Add(1)
    go func(feedURL string) {
      defer wg.Done()

      var resp *http.Response
      var bodyString string
      var fetchErr error

      for i := 0; i < maxRetries; i++ {
        resp, fetchErr = httpClient.Get(feedURL)
        if fetchErr == nil {
          bodyBytes := new(bytes.Buffer)
          bodyBytes.ReadFrom(resp.Body)
          bodyString = bodyBytes.String()
          resp.Body.Close()
          break
        }
        logError("Get RSS error: %s (attempt %d/%d): %v", feedURL, i+1, maxRetries, fetchErr)
        time.Sleep(retryInterval)
      }

      if fetchErr != nil {
        logError("Failed to fetch RSS: %s: %v", feedURL, fetchErr)
        return
      }

      cleanBody := cleanXMLContent(bodyString)

      var feed *gofeed.Feed
      var parseErr error
      for i := 0; i < maxRetries; i++ {
        feed, parseErr = fp.ParseString(cleanBody)
        if parseErr == nil {
          break
        }
        logError("Parse RSS error: %s (attempt %d/%d): %v", feedURL, i+1, maxRetries, parseErr)
        time.Sleep(retryInterval)
      }

      if parseErr != nil {
        logError("Failed to parse RSS: %s: %v", feedURL, parseErr)
        return
      }

      if len(feed.Items) == 0 {
        return
      }

      item := feed.Items[0]
      publishedTime := time.Now()
      if item.PublishedParsed != nil {
        publishedTime = *item.PublishedParsed
      } else if item.UpdatedParsed != nil {
        publishedTime = *item.UpdatedParsed
      }

      name := feed.Title
      nameMapping := map[string]string{
        "obaby@mars":                 "obaby",
        "青山小站 | 一个在帝都搬砖的新时代农民工": "青山小站",
        "Homepage on Miao Yu | 于淼":    "于淼",
        "Homepage on Yihui Xie | 谢益辉": "谢益辉",
      }
      if mapped, ok := nameMapping[name]; ok {
        name = mapped
      }

      avatarURL := avatars[name]
      if avatarURL == "" {
        avatarURL = "https://cos.lhasa.icu/LinksAvatar/default.png"
      }

      domainName := ""
      if feed.Link != "" {
        if domain, err := extractDomain(feed.Link); err == nil {
          domainName = domain
        }
      }

      mu.Lock()
      articles = append(articles, Article{
        DomainName: domainName,
        Name:       name,
        Title:      item.Title,
        Link:       item.Link,
        Avatar:     avatarURL,
        Date:       publishedTime.Format("2006-01-02"),
        DateTime:   publishedTime,
      })
      mu.Unlock()
    }(feedURL)
  }

  wg.Wait()

  sort.Slice(articles, func(i, j int) bool {
    return articles[i].DateTime.After(articles[j].DateTime)
  })

  return articles, nil
}

func writeJSON(path string, data []Article) error {
  if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
    return err
  }

  jsonData, err := json.MarshalIndent(data, "", "  ")
  if err != nil {
    return err
  }

  return os.WriteFile(path, jsonData, 0644)
}

func main() {
  feeds, err := readFeeds(feedsPath)
  if err != nil {
    logError("Read feeds error: %v", err)
    fmt.Printf("Error reading feeds: %v\n", err)
    return
  }

  avatars, err := loadAvatars(avatarsPath)
  if err != nil {
    logError("Load avatars error: %v", err)
    fmt.Printf("Error loading avatars: %v\n", err)
    return
  }

  articles, err := fetchRSS(feeds, avatars)
  if err != nil {
    logError("Fetch RSS error: %v", err)
    fmt.Printf("Error fetching RSS: %v\n", err)
    return
  }

  manual := Article{
    DomainName: "https://foreverblog.cn",
    Name:       "十年之约",
    Title:      "穿梭虫洞-随机访问十年之约友链博客",
    Link:       "https://foreverblog.cn/go.html",
    Date:       "2000-01-01",
    Avatar:     "https://cos.lhasa.icu/LinksAvatar/foreverblog.cn.png",
    DateTime:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
  }
  articles = append(articles, manual)

  if err := writeJSON(outputPath, articles); err != nil {
    logError("Write rss_data.json error: %v", err)
    fmt.Printf("Error writing rss_data.json: %v\n", err)
    return
  }

  fmt.Println("RSS fetch completed.")
}
