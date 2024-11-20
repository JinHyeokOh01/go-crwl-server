package crwl

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"
    "github.com/PuerkitoBio/goquery"
    "github.com/gin-gonic/gin"
)

type Notice struct {
    Title string
    Date  string
    Link  string
}

func GetSW(c *gin.Context) {
    url := "https://swedu.khu.ac.kr/bbs/board.php?bo_table=07_01"
    notices, err := crwlKHUNotices(url)
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

func crwlSWNotices(url string) ([]Notice, error) {
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
        notice := Notice{}

        // 제목과 링크
		titleLink := s.Find(".bo_tit a")
		notice.Title = strings.TrimSpace(titleLink.Text())

		// 링크 가져오기
		if link, exists := titleLink.Attr("href"); exists {
			notice.Link = link
		}

        // 날짜
        notice.Date = strings.TrimSpace(s.Find("td.td_datetime").Text())

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