package crwl

import (
    "net/http"
    "strings"
    "time"
    "sort"
    "github.com/PuerkitoBio/goquery"
    "github.com/gin-gonic/gin"
    "github.com/JinHyeokOh01/go-crwl-server/models"
)

const cseURL = "https://ce.khu.ac.kr/ce/user/bbs/BMSR00040/list.do?menuNo=1600045"

// 최적화된 버전
func GetCSE(c *gin.Context) {
    url := cseURL
    notices, err := crwlCSENotices(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    noticeRepo := repository.NewNoticeRepository()

    // 한 번의 쿼리로 모든 기존 공지사항 번호 가져오기
    existingNumbers, err := noticeRepo.GetAllNumbers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "데이터 확인 실패: " + err.Error()})
        return
    }

    // 맵을 사용하여 O(1) 검색
    existingMap := make(map[string]bool)
    for _, num := range existingNumbers {
        existingMap[num] = true
    }

    // 새로운 공지사항 필터링
    var newNotices []models.Notice
    for _, notice := range notices {
        if !existingMap[notice.Number] {
            newNotices = append(newNotices, notice)
        }
    }

    // 모든 데이터 저장
    err = noticeRepo.CreateBatch(notices)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "데이터 저장 실패: " + err.Error()})
        return
    }

    if len(newNotices) > 0 {
        c.JSON(http.StatusOK, gin.H{
            "message": "새로운 공지사항이 발견되었습니다",
            "count": len(newNotices),
            "notices": newNotices,
        })
    } else {
        c.JSON(http.StatusOK, gin.H{
            "message": "새로운 공지사항이 없습니다",
            "count": 0,
        })
    }
}

// repository/notice_repository.go에 추가
func (r *NoticeRepository) GetAllNumbers() ([]string, error) {
    rows, err := r.db.Query("SELECT number FROM notices")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var numbers []string
    for rows.Next() {
        var number string
        if err := rows.Scan(&number); err != nil {
            return nil, err
        }
        numbers = append(numbers, number)
    }
    return numbers, nil
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