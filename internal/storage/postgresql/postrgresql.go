package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"mmskazak/shorturl/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreSQL struct {
	db *sql.DB
}

func NewPostgreSQL(cfg *config.Config) (*PostgreSQL, error) {
	// Подключение к базе данных postgres для управления базами данных
	db, err := sql.Open("pgx", cfg.DataBaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to PostgreSQL: %w", err)
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

	// Проверяем наличие базы данных dbshorturl и создаем ее при необходимости
	var dbExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT datname FROM pg_catalog.pg_database WHERE datname = 'dbshorturl')").Scan(&dbExists) //nolint:lll //так удобнее читать
	if err != nil {
		return nil, fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !dbExists {
		_, err = db.Exec("CREATE DATABASE dbshorturl")
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	}

	// Закрываем соединение с базой данных postgres
	err = db.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close PostgreSQL connection: %w", err)
	}

	// Подключаемся к базе данных dbshorturl
	dbShortURL, err := sql.Open("pgx", cfg.DataBaseDSN+" dbname=dbshorturl")
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to dbshorturl: %w", err)
	}

	// Проверяем соединение с dbshorturl
	err = dbShortURL.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping dbshorturl connection: %w", err)
	}

	// Создаем таблицу shorturl, если она не существует
	_, err = dbShortURL.Exec(`
		CREATE TABLE IF NOT EXISTS shorturl (
			id SERIAL PRIMARY KEY,
			short_url VARCHAR(10) NOT NULL,
			original_url TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table shorturl: %w", err)
	}

	return &PostgreSQL{
		dbShortURL,
	}, nil
}

func (pg *PostgreSQL) GetShortURL(shortURL string) (string, error) {
	var originalURL string
	// Выполняем запрос SQL для получения original_url по short_url
	err := pg.db.QueryRow("SELECT original_url FROM shorturl WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("short URL not found %w", err)
		}
		return "", fmt.Errorf("failed to get original URL: %w", err)
	}
	return originalURL, nil
}

func (pg *PostgreSQL) SetShortURL(shortURL string, targetURL string) error {
	// Вставляем запись в базу данных
	_, err := pg.db.Exec("INSERT INTO shorturl (short_url, original_url) VALUES ($1, $2)", shortURL, targetURL)
	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	return nil
}

func (pg *PostgreSQL) Ping() error {
	err := pg.db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	return nil
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
