// main.go
package main

import(
    "log"
    "github.com/JinHyeokOh01/go-crwl-server/controllers"
    "github.com/JinHyeokOh01/go-crwl-server/db"
    "github.com/JinHyeokOh01/go-crwl-server/crwl"
    "github.com/gin-gonic/gin"
)

func main() {
    // 데이터베이스 초기화
    if err := db.Initialize(); err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    r := gin.Default()

    // 크롤링 엔드포인트 (자동으로 DB 저장)
    r.GET("/cse", crwl.GetCSE)
    r.GET("/sw", crwl.GetSW)

    // DB 조회용 API
    api := r.Group("/api")
    {
        api.GET("/notices", controllers.GetNotices)           // 전체 공지사항 조회
        api.GET("/notices/:number", controllers.GetNotice)    // 특정 공지사항 조회
    }

    r.Run(":5000")
}