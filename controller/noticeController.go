// controllers/notice_controller.go
package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/JinHyeokOh01/go-crwl-server/models"
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

// FetchAndSaveNotices CSE 공지사항 가져와서 저장
func (nc *NoticeController) FetchAndSaveNotices(c *gin.Context) {
    err := nc.service.FetchAndSaveNotices()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    notices, err := nc.service.GetAllNotices()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, notices)
}

// GetNotices 모든 공지사항 조회
func (nc *NoticeController) GetNotices(c *gin.Context) {
    notices, err := nc.service.GetAllNotices()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, notices)
}

// GetNotice 특정 공지사항 조회
func (nc *NoticeController) GetNotice(c *gin.Context) {
    number := c.Param("number")
    notice, err := nc.service.GetNotice(number)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, notice)
}

// CreateNotice 새 공지사항 생성
func (nc *NoticeController) CreateNotice(c *gin.Context) {
    var notice models.Notice
    if err := c.BindJSON(&notice); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := nc.service.CreateNotice(&notice)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, notice)
}

// UpdateNotice 공지사항 업데이트
func (nc *NoticeController) UpdateNotice(c *gin.Context) {
    number := c.Param("number")
    var notice models.Notice
    if err := c.BindJSON(&notice); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    notice.Number = number

    err := nc.service.UpdateNotice(&notice)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "공지사항이 업데이트되었습니다"})
}

// DeleteNotice 공지사항 삭제
func (nc *NoticeController) DeleteNotice(c *gin.Context) {
    number := c.Param("number")
    err := nc.service.DeleteNotice(number)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "공지사항이 삭제되었습니다"})
}