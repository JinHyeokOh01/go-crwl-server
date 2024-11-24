// repository/notice_repository.go
package repository

import (
    "database/sql"
    "github.com/JinHyeokOh01/go-crwl-server/models"
	"github.com/JinHyeokOh01/go-crwl-server/db"
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

// CreateBatchSW SW 공지사항 일괄 저장
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

// GetAllCSE CSE 공지사항 모두 조회
func (r *NoticeRepository) GetAllCSE() ([]models.Notice, error) {
    rows, err := r.db.Query("SELECT number, title, date, link FROM cse_notices ORDER BY date DESC")
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

// GetAllSW SW 공지사항 모두 조회
func (r *NoticeRepository) GetAllSW() ([]models.Notice, error) {
    rows, err := r.db.Query("SELECT number, title, date, link FROM sw_notices ORDER BY date DESC")
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

// GetCSENumbers CSE 공지사항 번호 목록 조회
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

// GetSWNumbers SW 공지사항 번호 목록 조회
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