// main.go
package main

import(
    "log"
    "github.com/JinHyeokOh01/go-crwl-server/controllers"
    "github.com/JinHyeokOh01/go-crwl-server/db"
    "github.com/JinHyeokOh01/go-crwl-server/crwl"
    "github.com/JinHyeokOh01/go-crwl-server/repository"
    "github.com/JinHyeokOh01/go-crwl-server/services"
    "github.com/gin-gonic/gin"
)

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

    r.Run(":5000")
}