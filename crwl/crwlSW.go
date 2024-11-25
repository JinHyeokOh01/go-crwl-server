package crwl

import (
    "net/http"
    "strings"
    "sort"
    "time"
    "github.com/PuerkitoBio/goquery"
    "github.com/gin-gonic/gin"
    "github.com/JinHyeokOh01/go-crwl-server/models"
    "github.com/JinHyeokOh01/go-crwl-server/repository"
    "github.com/JinHyeokOh01/go-crwl-server/services"

    "net/url"
)

func getIDFromURL(urlStr string) string {
    parsedURL, err := url.Parse(urlStr)
    if err != nil {
        return ""
    }
    values, _ := url.ParseQuery(parsedURL.RawQuery)
    return values.Get("wr_id")
}

func GetSW(c *gin.Context) {
    url := "https://swedu.khu.ac.kr/bbs/board.php?bo_table=07_01"
    crawledNotices, err := crwlSWNotices(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Service 계층 초기화
    noticeRepo := repository.NewNoticeRepository()
    noticeService := services.NewNoticeService(noticeRepo)

    // DB의 현재 공지사항 목록 조회
    dbNotices, err := noticeService.GetAllSWNotices()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "기존 데이터 조회 실패: " + err.Error(),
        })
        return
    }

    // 동기화가 필요한 항목들 찾기
    toAdd, toDelete := syncNotices(crawledNotices, dbNotices)

    // 삭제할 공지사항이 있다면 삭제
    if len(toDelete) > 0 {
        err = noticeService.DeleteBatchSW(toDelete)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "삭제 실패: " + err.Error(),
            })
            return
        }
    }

    // 추가할 공지사항이 있다면 추가
    if len(toAdd) > 0 {
        err = noticeService.CreateBatchSW(toAdd)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "저장 실패: " + err.Error(),
            })
            return
        }
    }

    // 응답 반환
    if len(toAdd) > 0 {
        c.JSON(http.StatusOK, gin.H{
            "message": "새로운 소프트웨어중심대학사업단 공지사항이 있습니다",
            "notices": toAdd,
            "sync_status": gin.H{
                "added": len(toAdd),
                "deleted": len(toDelete),
            },
        })
    } else {
        c.JSON(http.StatusOK, gin.H{
            "message": "새로운 소프트웨어중심대학사업단 공지사항이 없습니다",
            "sync_status": gin.H{
                "added": 0,
                "deleted": len(toDelete),
            },
        })
    }
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

    // HTML에서 해당 태그와 일치하는 부분에서 데이터 가져오기
    // 컴공과와 소중단의 HTML 페이지 구조는 약간 다름.
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
    // 날짜순 정렬
    sort.Sort(NoticeSlice(notices))

    return notices, nil
}