package crwl

import (
    "net/http"
    "strings"
    "time"
    "sort"
    "github.com/PuerkitoBio/goquery"
    "github.com/gin-gonic/gin"
    "github.com/JinHyeokOh01/go-crwl-server/models"
    "github.com/JinHyeokOh01/go-crwl-server/repository"
    "github.com/JinHyeokOh01/go-crwl-server/services"
)

const cseURL = "https://ce.khu.ac.kr/ce/user/bbs/BMSR00040/list.do?menuNo=1600045"

func GetCSE(c *gin.Context) {
    url := cseURL
    crawledNotices, err := crwlCSENotices(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Service 계층 초기화
    noticeRepo := repository.NewNoticeRepository()
    noticeService := services.NewNoticeService(noticeRepo)
    
    // DB에서 현재 저장된 모든 공지사항 조회
    dbNotices, err := noticeService.GetAllCSENotices()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "DB 조회 실패: " + err.Error(),
        })
        return
    }

    // 크롤링된 데이터와 DB 데이터를 비교하여 동기화
    toAdd, toDelete := syncNotices(crawledNotices, dbNotices)

    // 새로운 공지사항 추가
    if len(toAdd) > 0 {
        err = noticeService.CreateBatchCSE(toAdd)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "새 공지사항 저장 실패: " + err.Error(),
            })
            return
        }
    }

    // 삭제된 공지사항 제거
    if len(toDelete) > 0 {
        err = noticeService.DeleteBatchCSE(toDelete)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "공지사항 삭제 실패: " + err.Error(),
            })
            return
        }
    }

    // 응답 반환
    if len(toAdd) > 0 {
        c.JSON(http.StatusOK, gin.H{
            "message": "새로운 컴퓨터공학과 공지사항이 있습니다",
            "notices": toAdd,
            "sync_status": gin.H{
                "added": len(toAdd),
                "deleted": len(toDelete),
            },
        })
    } else {
        c.JSON(http.StatusOK, gin.H{
            "message": "새로운 컴퓨터공학과 공지사항이 없습니다",
            "sync_status": gin.H{
                "added": 0,
                "deleted": len(toDelete),
            },
        })
    }
}

func syncNotices(crawled []models.Notice, dbNotices []models.Notice) ([]models.Notice, []models.Notice) {
    crawledMap := make(map[string]models.Notice)
    dbMap := make(map[string]models.Notice)

    // 맵으로 변환하여 빠른 검색 가능하게 함
    for _, notice := range crawled {
        crawledMap[notice.Number] = notice
    }
    for _, notice := range dbNotices {
        dbMap[notice.Number] = notice
    }

    // 추가할 항목 찾기 (크롤링된 데이터에는 있지만 DB에는 없는 항목)
    var toAdd []models.Notice
    for number, notice := range crawledMap {
        if _, exists := dbMap[number]; !exists {
            toAdd = append(toAdd, notice)
        }
    }

    // 삭제할 항목 찾기 (DB에는 있지만 크롤링된 데이터에는 없는 항목)
    var toDelete []models.Notice
    for number, notice := range dbMap {
        if _, exists := crawledMap[number]; !exists {
            toDelete = append(toDelete, notice)
        }
    }

    return toAdd, toDelete
}

func crwlCSENotices(url string) ([]models.Notice, error) {
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
    doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
        // 글 번호가 '대학'과 '공지'인 것은 딱히 필요가 없는 것 같고 파싱도 잘 안되서 버림.
        number := strings.TrimSpace(s.Find("td.align-middle").Text())
        if number == "대학" || number == "공지"{
            return
        }
        notice := models.Notice{}
        // 글 번호
        notice.Number = number
        // 제목
        titleLink := s.Find("td.tal a")
        notice.Title = strings.Join(strings.Fields(titleLink.Text()), " ")
        // 날짜
        notice.Date = strings.TrimSpace(s.Find("td:nth-child(4)").Text())
        // CSE 홈페이지는 개별 공지 글마다 URL이 존재하지 않음. Default URL 적용
        notice.Link = cseURL

        notices = append(notices, notice)
    })
    // 날짜순 정렬
    sort.Sort(NoticeSlice(notices))
    return notices, nil
}