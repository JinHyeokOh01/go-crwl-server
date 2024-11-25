package controllers

import (
   "net/http"
   "github.com/gin-gonic/gin"
   "github.com/JinHyeokOh01/go-crwl-server/services"
)

type NoticeController struct {
   service *services.NoticeService
}

func NewNoticeController(service *services.NoticeService) *NoticeController {
   return &NoticeController{
       service: service,
   }
}

// GetCSENotices CSE 공지사항 조회
func (nc *NoticeController) GetCSENotices(c *gin.Context) {
   notices, err := nc.service.GetAllCSENotices()
   if err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   c.JSON(http.StatusOK, notices)
}

// GetSWNotices SW 공지사항 조회
func (nc *NoticeController) GetSWNotices(c *gin.Context) {
   notices, err := nc.service.GetAllSWNotices()
   if err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   c.JSON(http.StatusOK, notices)
}

// DeleteAllNotices 모든 공지사항 삭제
func (nc *NoticeController) DeleteAllNotices(c *gin.Context) {
   // CSE 공지사항 삭제
   if err := nc.service.DeleteAllCSE(); err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{
           "error": "CSE 공지사항 삭제 실패: " + err.Error(),
       })
       return
   }

   // SW 공지사항 삭제
   if err := nc.service.DeleteAllSW(); err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{
           "error": "SW 공지사항 삭제 실패: " + err.Error(),
       })
       return
   }

   c.JSON(http.StatusOK, gin.H{
       "message": "모든 공지사항이 삭제되었습니다",
   })
}

// DeleteAllCSENotices CSE 공지사항만 삭제
func (nc *NoticeController) DeleteAllCSENotices(c *gin.Context) {
   if err := nc.service.DeleteAllCSE(); err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{
           "error": "CSE 공지사항 삭제 실패: " + err.Error(),
       })
       return
   }

   c.JSON(http.StatusOK, gin.H{
       "message": "모든 CSE 공지사항이 삭제되었습니다",
   })
}

// DeleteAllSWNotices SW 공지사항만 삭제
func (nc *NoticeController) DeleteAllSWNotices(c *gin.Context) {
   if err := nc.service.DeleteAllSW(); err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{
           "error": "SW 공지사항 삭제 실패: " + err.Error(),
       })
       return
   }

   c.JSON(http.StatusOK, gin.H{
       "message": "모든 SW 공지사항이 삭제되었습니다",
   })
}