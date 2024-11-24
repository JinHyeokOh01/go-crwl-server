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