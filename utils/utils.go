package utils

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ConnectToDB(connStr string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	return db, nil
}
