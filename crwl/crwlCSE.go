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
)

const cseURL = "https://ce.khu.ac.kr/ce/user/bbs/BMSR00040/list.do?menuNo=1600045"

func GetCSE(c *gin.Context) {
    url := cseURL
    notices, err := crwlCSENotices(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 데이터베이스에 자동 저장
    noticeRepo := repository.NewNoticeRepository()

    // 기존 공지사항 번호들 조회
    existingNotices, err := noticeRepo.GetCSENumbers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "기존 데이터 조회 실패: " + err.Error()})
        return
    }

    // 새로운 공지사항 필터링
    existingMap := make(map[string]bool)
    for _, num := range existingNotices {
        existingMap[num] = true
    }

    var newNotices []models.Notice
    for _, notice := range notices {
        if !existingMap[notice.Number] {
            newNotices = append(newNotices, notice)
        }
    }

    // DB에 저장
    err = noticeRepo.CreateBatchCSE(notices)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "저장 실패: " + err.Error()})
        return
    }

    // 새로운 공지사항만 응답으로 반환
    if len(newNotices) > 0 {
        c.JSON(http.StatusOK, gin.H{
            "message": "새로운 CSE 공지사항이 발견되었습니다",
            "count": len(newNotices),
            "notices": newNotices,
        })
    } else {
        c.JSON(http.StatusOK, gin.H{
            "message": "새로운 CSE 공지사항이 없습니다",
            "count": 0,
        })
    }
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

    //HTML에서 해당 태그와 일치하는 부분에서 데이터 가져오기
    doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
        //글 번호가 '대학'과 '공지'인 것은 딱히 필요가 없는 것 같고 파싱도 잘 안되서 버림.
        number := strings.TrimSpace(s.Find("td.align-middle").Text())
        if number == "대학" || number == "공지"{
            return
        }
        notice := models.Notice{}
        //글 번호
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
    //날짜순 정렬
    sort.Sort(NoticeSlice(notices))
    return notices, nil
}