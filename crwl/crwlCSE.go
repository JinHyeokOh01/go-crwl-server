package crwl

import (
    "fmt"
    "net/http"
    "strings"
    "time"
    "github.com/PuerkitoBio/goquery"
    "github.com/gin-gonic/gin"
)

type CSENotice struct {
    Title string
    Date  string
}

func GetCSE(c *gin.Context) {
    url := "https://ce.khu.ac.kr/ce/user/bbs/BMSR00040/list.do?menuNo=1600045"
    notices, err := crwlCSENotices(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, notices)

    /*
    for _, notice := range notices {
        fmt.Printf("제목: %s\n날짜: %s\n링크: %s\n---\n", notice.Title, notice.Date, notice.Link)
    }
    */
}

func crwlCSENotices(url string) ([]CSENotice, error) {
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

    var notices []CSENotice

    doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
        number := strings.TrimSpace(s.Find("td.align-middle").Text())
        if number == "대학" || number == "공지" {
            return
        }

        notice := CSENotice{}

        // 제목과 링크
        titleLink := s.Find("td.tal a")
        notice.Title = strings.TrimSpace(titleLink.Text())

        // 날짜
        notice.Date = strings.TrimSpace(s.Find("td:nth-child(4)").Text())

        notices = append(notices, notice)
    })

    return notices, nil
}

func crwlAllPages(baseURL string, maxPages int) ([]CSENotice, error) {
    var allNotices []CSENotice

    for page := 1; page <= maxPages; page++ {
        pageURL := fmt.Sprintf("%s&pageIndex=%d", baseURL, page)
        notices, err := crwlCSENotices(pageURL)
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