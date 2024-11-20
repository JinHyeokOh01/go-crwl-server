package crwl

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"
    "github.com/PuerkitoBio/goquery"
)

type Notice struct {
    Title string
    Date  string
    Link  string
}

func GetCSE() {
    url := "https://ce.khu.ac.kr/ce/user/bbs/BMSR00040/list.do?menuNo=1600045"
    notices, err := crawlCSENotices(url)
    if err != nil {
        log.Fatal(err)
    }

    // 결과 출력
    for _, notice := range notices {
        fmt.Printf("제목: %s\n날짜: %s\n링크: %s\n---\n", notice.Title, notice.Date, notice.Link)
    }
}

func crawlCSENotices(url string) ([]Notice, error) {
    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, err
    }

    var notices []Notice

    doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
        number := strings.TrimSpace(s.Find("td.align-middle").Text())
        if number == "대학" || number == "공지" {
            return
        }

        notice := Notice{}

        // 제목과 링크
        titleLink := s.Find("td.tal a")
        notice.Title = strings.TrimSpace(titleLink.Text())
        
        // 링크 생성
        href, _ := titleLink.Attr("href")
        if href != "" {
            idStart := strings.Index(href, "'") + 1
            idEnd := strings.LastIndex(href, "'")
            if idStart > 0 && idEnd > idStart {
                noticeID := href[idStart:idEnd]
                notice.Link = fmt.Sprintf("https://ce.khu.ac.kr/ce/user/bbs/BMSR00040/view.do?nttId=%s&menuNo=1600045", noticeID)
            }
        }

        // 날짜
        notice.Date = strings.TrimSpace(s.Find("td:nth-child(4)").Text())

        notices = append(notices, notice)
    })

    return notices, nil
}

func crawlAllPages(baseURL string, maxPages int) ([]Notice, error) {
    var allNotices []Notice

    for page := 1; page <= maxPages; page++ {
        pageURL := fmt.Sprintf("%s&pageIndex=%d", baseURL, page)
        notices, err := crawlKHUNotices(pageURL)
        if err != nil {
            return nil, err
        }

        if len(notices) == 0 {
            break
        }

        allNotices = append(allNotices, notices...)
        time.Sleep(1 * time.Second)
    }

    return allNotices, nil
}