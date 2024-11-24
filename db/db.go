package db

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Initialize() error {
    dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
        "mydb",           // MYSQL_USER
        "mydb",           // MYSQL_PASSWORD
        "db",            // MYSQL_HOST
        "mydb",          // MYSQL_DATABASE
    )

    var err error
    DB, err = sql.Open("mysql", dataSourceName)
    if err != nil {
        return fmt.Errorf("데이터베이스 연결 실패: %v", err)
    }

    err = createTables()
    if err != nil {
        return fmt.Errorf("테이블 생성 실패: %v", err)
    }

    return nil
}

func createTables() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS cse_notices (
            number VARCHAR(255) PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            date VARCHAR(255) NOT NULL,
            link VARCHAR(255) NOT NULL
        )`,
        `CREATE TABLE IF NOT EXISTS sw_notices (
            number VARCHAR(255) PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            date VARCHAR(255) NOT NULL,
            link VARCHAR(255) NOT NULL
        )`,
    }

    for _, query := range queries {
        _, err := DB.Exec(query)
        if err != nil {
            return fmt.Errorf("테이블 생성 실패: %v", err)
        }
    }
    return nil
}

func Close() error {
    if DB != nil {
        return DB.Close()
    }
    return nil
}