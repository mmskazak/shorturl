package postgresql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreSQL struct {
	db *sql.DB
}

func NewPostgreSQL(connectionString string) (*PostgreSQL, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to postgresql: %w", err)
	}

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		err := db.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close PostgreSQL connection after db.Ping: %w", err)
		}
		return nil, fmt.Errorf("failed to ping PostgreSQL connection: %w", err)
	}

	return &PostgreSQL{
		db: db,
	}, nil
}

func (pg *PostgreSQL) Close() {
	if pg.db != nil {
		err := pg.db.Close()
		if err != nil {
			log.Fatalf("Error closing database connection: %v\n", err)
		}
		log.Println("Database connection closed.")
	}
}

func (pg *PostgreSQL) GetShortURL(id string) (string, error) {
	return id, nil
}

func (pg *PostgreSQL) SetShortURL(id string, targetURL string) error {
	return nil
}

func (pg *PostgreSQL) Ping() error {
	err := pg.db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	return nil
}
