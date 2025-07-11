package internal

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"log"
)

type Database struct {
	DB *sqlx.DB
}

func NewDatabaseConnection(dbDriverStr string, connectionStr string) (*Database, error) {
	database, err := sqlx.Connect(dbDriverStr, connectionStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	err = database.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка пинга БД: %w", err)
	}

	log.Println("подключение к БД успешно установлено")
	return &Database{DB: database}, nil
}

func (database *Database) Close() error {
	err := database.DB.Close()
	if err != nil {
		return fmt.Errorf("ошибка закрытия подключения к БД: %v", err)
	}
	return nil
}
