package crwl

import (
    "net/http"
    "strings"
    "sort"
    "time"
    "github.com/PuerkitoBio/goquery"
    "github.com/gin-gonic/gin"
    "github.com/JinHyeokOh01/go-crwl-server/models"

    "net/url"
)

func getIDFromURL(urlStr string) string {
    // 1. URL을 파싱합니다
    parsedURL, err := url.Parse(urlStr)
    if err != nil {
        return ""
    }
    
    // 2. 쿼리 파라미터를 파싱합니다
    values, _ := url.ParseQuery(parsedURL.RawQuery)
    
    // 3. wr_id 값을 가져옵니다
    return values.Get("wr_id")
}

func GetSW(c *gin.Context) {
    url := "https://swedu.khu.ac.kr/bbs/board.php?bo_table=07_01"
    notices, err := crwlSWNotices(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    //결과값을 JSON으로 반환
    c.JSON(http.StatusOK, notices)
}

func crwlSWNotices(url string) ([]models.Notice, error) {
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

    var notices []models.Notice

    //HTML에서 해당 태그와 일치하는 부분에서 데이터 가져오기
    //컴공과와 소중단의 HTML 페이지 구조는 약간 다름.
    doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
        notice := models.Notice{}
        // 제목
		titleLink := s.Find(".bo_tit a")
		notice.Title = strings.TrimSpace(titleLink.Text())

		// 링크 가져오기
		if link, exists := titleLink.Attr("href"); exists {
			notice.Link = link
            //PRIMARY KEY로 사용하기 위해 글의 고유 ID를 가져옴
            notice.Number = getIDFromURL(link)
		}

        // 날짜
        notice.Date = strings.TrimSpace(s.Find("td.td_datetime").Text())

        notices = append(notices, notice)
    })
    //날짜순 정렬
    sort.Sort(NoticeSlice(notices))

    return notices, nil
}