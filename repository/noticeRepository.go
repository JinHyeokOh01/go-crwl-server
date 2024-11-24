// repository/notice_repository.go
package repository

import (
    "database/sql"
    "fmt"
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

// Create 새 공지사항 생성
func (r *NoticeRepository) Create(notice *models.Notice) error {
    query := `
        INSERT INTO notices (number, title, date, link)
        VALUES (?, ?, ?, ?)
    `
    result, err := r.db.Exec(query, notice.Number, notice.Title, notice.Date, notice.Link)
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    notice.ID = id
    return nil
}

// GetAll 모든 공지사항 조회
func (r *NoticeRepository) GetAll() ([]models.Notice, error) {
    query := `
        SELECT id, number, title, date, link, created_at, updated_at 
        FROM notices 
        ORDER BY date DESC
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notices []models.Notice
    for rows.Next() {
        var n models.Notice
        err := rows.Scan(&n.ID, &n.Number, &n.Title, &n.Date, &n.Link, &n.CreatedAt, &n.UpdatedAt)
        if err != nil {
            return nil, err
        }
        notices = append(notices, n)
    }
    return notices, nil
}

// GetByNumber 특정 번호의 공지사항 조회
func (r *NoticeRepository) GetByNumber(number string) (*models.Notice, error) {
    query := `
        SELECT id, number, title, date, link, created_at, updated_at 
        FROM notices 
        WHERE number = ?
    `
    var notice models.Notice
    err := r.db.QueryRow(query, number).Scan(
        &notice.ID, &notice.Number, &notice.Title, &notice.Date, 
        &notice.Link, &notice.CreatedAt, &notice.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("공지사항을 찾을 수 없습니다")
    }
    if err != nil {
        return nil, err
    }
    return &notice, nil
}

// Update 공지사항 업데이트
func (r *NoticeRepository) Update(notice *models.Notice) error {
    query := `
        UPDATE notices 
        SET title = ?, date = ?, link = ? 
        WHERE number = ?
    `
    result, err := r.db.Exec(query, notice.Title, notice.Date, notice.Link, notice.Number)
    if err != nil {
        return err
    }

    affected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if affected == 0 {
        return fmt.Errorf("공지사항을 찾을 수 없습니다")
    }
    return nil
}

// Delete 공지사항 삭제
func (r *NoticeRepository) Delete(number string) error {
    query := `DELETE FROM notices WHERE number = ?`
    result, err := r.db.Exec(query, number)
    if err != nil {
        return err
    }

    affected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if affected == 0 {
        return fmt.Errorf("공지사항을 찾을 수 없습니다")
    }
    return nil
}

// CreateBatch 여러 공지사항 한 번에 생성/업데이트
func (r *NoticeRepository) CreateBatch(notices []models.Notice) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare(`
        INSERT INTO notices (number, title, date, link)
        VALUES (?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            title = VALUES(title),
            date = VALUES(date),
            link = VALUES(link),
            updated_at = CURRENT_TIMESTAMP
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