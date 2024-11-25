package repository

import (
    "database/sql"
    "github.com/JinHyeokOh01/go-crwl-server/models"
    "github.com/JinHyeokOh01/go-crwl-server/db"
    "strings"
)

type NoticeRepository struct {
    db *sql.DB
}

func NewNoticeRepository() *NoticeRepository {
    return &NoticeRepository{
        db: db.DB,
    }
}

// CreateBatchCSE CSE 공지사항 일괄 저장
func (r *NoticeRepository) CreateBatchCSE(notices []models.Notice) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare(`
        INSERT INTO cse_notices (number, title, date, link)
        VALUES (?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            title = VALUES(title),
            date = VALUES(date),
            link = VALUES(link)
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, notice := range notices {
        _, err = stmt.Exec(notice.Number, notice.Title, notice.Date, notice.Link)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

// DeleteBatchCSE CSE 공지사항 일괄 삭제
func (r *NoticeRepository) DeleteBatchCSE(notices []models.Notice) error {
    if len(notices) == 0 {
        return nil
    }

    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // numbers를 IN 절에서 사용하기 위한 플레이스홀더 생성
    placeholders := make([]string, len(notices))
    args := make([]interface{}, len(notices))
    
    for i, notice := range notices {
        placeholders[i] = "?"
        args[i] = notice.Number
    }

    // DELETE 쿼리 실행
    query := "DELETE FROM cse_notices WHERE number IN (" + 
        strings.Join(placeholders, ",") + ")"
    
    _, err = tx.Exec(query, args...)
    if err != nil {
        return err
    }

    return tx.Commit()
}

// GetAllCSENotices CSE 공지사항 전체 조회 (이름 변경 및 정렬 조건 추가)
func (r *NoticeRepository) GetAllCSENotices() ([]models.Notice, error) {
    rows, err := r.db.Query(`
        SELECT number, title, date, link 
        FROM cse_notices 
        ORDER BY date DESC, number DESC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notices []models.Notice
    for rows.Next() {
        var n models.Notice
        if err := rows.Scan(&n.Number, &n.Title, &n.Date, &n.Link); err != nil {
            return nil, err
        }
        notices = append(notices, n)
    }
    return notices, nil
}

// CreateBatchSW SW 공지사항 일괄 저장 (SW 관련 함수들도 동일한 패턴으로 수정)
func (r *NoticeRepository) CreateBatchSW(notices []models.Notice) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare(`
        INSERT INTO sw_notices (number, title, date, link)
        VALUES (?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            title = VALUES(title),
            date = VALUES(date),
            link = VALUES(link)
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, notice := range notices {
        _, err = stmt.Exec(notice.Number, notice.Title, notice.Date, notice.Link)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

// DeleteBatchSW SW 공지사항 일괄 삭제
func (r *NoticeRepository) DeleteBatchSW(notices []models.Notice) error {
    if len(notices) == 0 {
        return nil
    }

    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    placeholders := make([]string, len(notices))
    args := make([]interface{}, len(notices))
    
    for i, notice := range notices {
        placeholders[i] = "?"
        args[i] = notice.Number
    }

    query := "DELETE FROM sw_notices WHERE number IN (" + 
        strings.Join(placeholders, ",") + ")"
    
    _, err = tx.Exec(query, args...)
    if err != nil {
        return err
    }

    return tx.Commit()
}

// GetAllSWNotices SW 공지사항 전체 조회 (이름 변경 및 정렬 조건 추가)
func (r *NoticeRepository) GetAllSWNotices() ([]models.Notice, error) {
    rows, err := r.db.Query(`
        SELECT number, title, date, link 
        FROM sw_notices 
        ORDER BY date DESC, number DESC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notices []models.Notice
    for rows.Next() {
        var n models.Notice
        if err := rows.Scan(&n.Number, &n.Title, &n.Date, &n.Link); err != nil {
            return nil, err
        }
        notices = append(notices, n)
    }
    return notices, nil
}

// GetCSENumbers와 GetSWNumbers는 이제 불필요할 수 있지만 호환성을 위해 유지
func (r *NoticeRepository) GetCSENumbers() ([]string, error) {
    rows, err := r.db.Query("SELECT number FROM cse_notices")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var numbers []string
    for rows.Next() {
        var number string
        if err := rows.Scan(&number); err != nil {
            return nil, err
        }
        numbers = append(numbers, number)
    }
    return numbers, nil
}

func (r *NoticeRepository) GetSWNumbers() ([]string, error) {
    rows, err := r.db.Query("SELECT number FROM sw_notices")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var numbers []string
    for rows.Next() {
        var number string
        if err := rows.Scan(&number); err != nil {
            return nil, err
        }
        numbers = append(numbers, number)
    }
    return numbers, nil
}

// DeleteAllCSE CSE 공지사항 전체 삭제
func (r *NoticeRepository) DeleteAllCSE() error {
    _, err := r.db.Exec("DELETE FROM cse_notices")
    return err
}

// DeleteAllSW SW 공지사항 전체 삭제
func (r *NoticeRepository) DeleteAllSW() error {
    _, err := r.db.Exec("DELETE FROM sw_notices")
    return err
}

