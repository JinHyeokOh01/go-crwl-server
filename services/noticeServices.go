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
    return s.repo.GetAllCSENotices()  // 메서드 이름 수정
}

// GetAllSWNotices SW 공지사항 모두 조회
func (s *NoticeService) GetAllSWNotices() ([]models.Notice, error) {
    return s.repo.GetAllSWNotices()  // 메서드 이름 수정
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

// DeleteBatchCSE CSE 공지사항 일괄 삭제 (추가)
func (s *NoticeService) DeleteBatchCSE(notices []models.Notice) error {
    return s.repo.DeleteBatchCSE(notices)
}

// DeleteBatchSW SW 공지사항 일괄 삭제 (추가)
func (s *NoticeService) DeleteBatchSW(notices []models.Notice) error {
    return s.repo.DeleteBatchSW(notices)
}

// DeleteAllCSE CSE 공지사항 전체 삭제
func (s *NoticeService) DeleteAllCSE() error {
    return s.repo.DeleteAllCSE()
}

// DeleteAllSW SW 공지사항 전체 삭제
func (s *NoticeService) DeleteAllSW() error {
    return s.repo.DeleteAllSW()
}