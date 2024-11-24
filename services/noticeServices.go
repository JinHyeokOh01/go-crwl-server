// services/notice_service.go
package services

import (
    "github.com/JinHyeokOh01/go-crwl-server/models"
    "github.com/JinHyeokOh01/go-crwl-server/repository"
)

type NoticeService struct {
    repo *repository.NoticeRepository
}

func NewNoticeService(repo *repository.NoticeRepository) *NoticeService {
    return &NoticeService{
        repo: repo,
    }
}

// GetAllCSENotices CSE 공지사항 모두 조회
func (s *NoticeService) GetAllCSENotices() ([]models.Notice, error) {
    return s.repo.GetAllCSE()
}

// GetAllSWNotices SW 공지사항 모두 조회
func (s *NoticeService) GetAllSWNotices() ([]models.Notice, error) {
    return s.repo.GetAllSW()
}

// GetCSENumbers CSE 공지사항 번호 목록 조회
func (s *NoticeService) GetCSENumbers() ([]string, error) {
    return s.repo.GetCSENumbers()
}

// GetSWNumbers SW 공지사항 번호 목록 조회
func (s *NoticeService) GetSWNumbers() ([]string, error) {
    return s.repo.GetSWNumbers()
}

// CreateBatchCSE CSE 공지사항 일괄 저장
func (s *NoticeService) CreateBatchCSE(notices []models.Notice) error {
    return s.repo.CreateBatchCSE(notices)
}

// CreateBatchSW SW 공지사항 일괄 저장
func (s *NoticeService) CreateBatchSW(notices []models.Notice) error {
    return s.repo.CreateBatchSW(notices)
}