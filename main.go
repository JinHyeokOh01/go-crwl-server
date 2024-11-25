package main

import(
    "log"
    "time"
    "net/http"
    "github.com/JinHyeokOh01/go-crwl-server/controllers"
    "github.com/JinHyeokOh01/go-crwl-server/db"
    "github.com/JinHyeokOh01/go-crwl-server/crwl"
    "github.com/JinHyeokOh01/go-crwl-server/repository"
    "github.com/JinHyeokOh01/go-crwl-server/services"
    "github.com/gin-gonic/gin"
)

func performCrawling() {
    log.Println("크롤링 시작...")
    
    endpoints := []string{"cse", "sw"}
    for _, endpoint := range endpoints {
        resp, err := http.Get("http://localhost:5000/" + endpoint)
        if err != nil {
            log.Printf("%s 크롤링 실패: %v\n", endpoint, err)
            continue
        }
        resp.Body.Close()
        log.Printf("%s 크롤링 완료\n", endpoint)
    }
 }

// 주기적 크롤링을 위한 함수
func startPeriodicCrawling() {
    // 즉시 한 번 실행
    performCrawling()
 
    // 이후 주기적 실행
    ticker := time.NewTicker(1 * time.Hour)
    go func() {
        for range ticker.C {
            performCrawling()
        }
    }()
 }
 
func main() {
    // 데이터베이스 초기화
    if err := db.Initialize(); err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 의존성 초기화
    noticeRepo := repository.NewNoticeRepository()
    noticeService := services.NewNoticeService(noticeRepo)
    noticeController := controllers.NewNoticeController(noticeService)

    r := gin.Default()

    // 크롤링 엔드포인트 (자동으로 DB 저장)
    r.GET("/cse", crwl.GetCSE)
    r.GET("/sw", crwl.GetSW)

    // DB 조회용 엔드포인트
    r.GET("/notices/cse", noticeController.GetCSENotices)
    r.GET("/notices/sw", noticeController.GetSWNotices)

    // DB 일괄 삭제 엔드포인트
    r.DELETE("/notices", noticeController.DeleteAllNotices)        // 전체 삭제
    r.DELETE("/notices/cse", noticeController.DeleteAllCSENotices) // CSE만 삭제
    r.DELETE("/notices/sw", noticeController.DeleteAllSWNotices)   // SW만 삭제

    // 서버 시작 후 주기적 크롤링 시작
    go func() {
        // 서버가 완전히 시작될 때까지 잠시 대기
        time.Sleep(20 * time.Second)
        startPeriodicCrawling()
    }()

    r.Run(":5000")
}