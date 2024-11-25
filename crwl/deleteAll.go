package crwl

import(
    "github.com/JinHyeokOh01/go-crwl-server/repository"
    "github.com/JinHyeokOh01/go-crwl-server/services"
    "github.com/gin-gonic/gin"

    "net/http"
)

func DeleteAllCSENotices(c *gin.Context) {
    noticeRepo := repository.NewNoticeRepository()
    noticeService := services.NewNoticeService(noticeRepo)

    if err := noticeService.DeleteAllCSE(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "컴퓨터공학과 공지사항 삭제 실패: " + err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "모든 컴퓨터공학과 공지사항이 삭제되었습니다",
    })
}

func DeleteAllSWNotices(c *gin.Context) {
    noticeRepo := repository.NewNoticeRepository()
    noticeService := services.NewNoticeService(noticeRepo)

    if err := noticeService.DeleteAllSW(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "소프트웨어중심대학사업단 공지사항 삭제 실패: " + err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "모든 소프트웨어중심대학사업단 공지사항이 삭제되었습니다",
    })
}