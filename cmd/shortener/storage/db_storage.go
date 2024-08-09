package storage

import (
	"database/sql"
	"fmt"
	"github.com/fngoc/url-shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStore struct {
	db *sql.DB
}

var postgresInstant DBStore

func InitializeDB(dbConf string) error {
	pqx, err := sql.Open("pgx", dbConf)
	if err != nil {
		return err
	}

	postgresInstant.db = pqx

	if err := createTables(pqx); err != nil {
		return err
	}
	Store = postgresInstant
	return nil
}

func (dbs DBStore) GetData(key string) (string, error) {
	row := dbs.db.QueryRow("SELECT original_url FROM url_shortener WHERE short_url = $1", key)
	var originalUrl string
	if err := row.Scan(&originalUrl); err != nil {
		return "", err
	}
	return originalUrl, nil
}

func (dbs DBStore) SaveData(key string, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("key or value is empty")
	}
	_, err := dbs.db.Exec("INSERT INTO url_shortener(short_url, original_url) VALUES ($1, $2)", key, value)
	if err != nil {
		return err
	}
	return nil
}

func CustomPing() bool {
	if postgresInstant.db == nil {
		return false
	}
	err := postgresInstant.db.Ping()
	if err != nil {
		return false
	}
	return true
}

func createTables(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS url_shortener (
		uuid SERIAL PRIMARY KEY,
		short_url VARCHAR NOT NULL UNIQUE,
		original_url VARCHAR NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	logger.Log.Info("Database table created")
	return nil
}
