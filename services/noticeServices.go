// services/notice_service.go
package services

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
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

// FetchAndSaveNotices CSE API에서 공지사항을 가져와서 저장
func (s *NoticeService) FetchAndSaveNotices() error {
    notices, err := s.fetchFromAPI()
    if err != nil {
        return err
    }

    return s.repo.CreateBatch(notices)
}

// fetchFromAPI CSE API에서 공지사항 가져오기
func (s *NoticeService) fetchFromAPI() ([]models.Notice, error) {
    resp, err := http.Get("http://localhost:5000/cse")
    if err != nil {
        return nil, fmt.Errorf("API 호출 실패: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("응답 읽기 실패: %v", err)
    }

    var notices []models.Notice
    if err := json.Unmarshal(body, &notices); err != nil {
        return nil, fmt.Errorf("JSON 파싱 실패: %v", err)
    }

    return notices, nil
}

// GetAllNotices 모든 공지사항 조회
func (s *NoticeService) GetAllNotices() ([]models.Notice, error) {
    return s.repo.GetAll()
}

// GetNotice 특정 공지사항 조회
func (s *NoticeService) GetNotice(number string) (*models.Notice, error) {
    return s.repo.GetByNumber(number)
}

// CreateNotice 새 공지사항 생성
func (s *NoticeService) CreateNotice(notice *models.Notice) error {
    return s.repo.Create(notice)
}

// UpdateNotice 공지사항 업데이트
func (s *NoticeService) UpdateNotice(notice *models.Notice) error {
    return s.repo.Update(notice)
}

// DeleteNotice 공지사항 삭제
func (s *NoticeService) DeleteNotice(number string) error {
    return s.repo.Delete(number)
}