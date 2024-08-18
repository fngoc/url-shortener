package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/fngoc/url-shortener/internal/logger"
	"github.com/fngoc/url-shortener/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type DBStore struct {
	db *sql.DB
}

type DBError struct {
	ShortURL string
	Err      *pgconn.PgError
}

func (p *DBError) Error() string {
	return p.Err.Message
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

func (dbs DBStore) GetData(ctx context.Context, key string) (string, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := dbs.db.QueryRowContext(dbCtx, "SELECT original_url FROM url_shortener WHERE short_url = $1", key)
	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		return "", err
	}
	return originalURL, nil
}

func (dbs DBStore) GetAllData(ctx context.Context) ([]models.ResponseDto, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := dbs.db.QueryContext(dbCtx, "SELECT short_url, original_url FROM url_shortener")
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var result []models.ResponseDto
	for rows.Next() {
		item := models.ResponseDto{}

		if err := rows.Scan(&item.ShortURL, &item.OriginalURL); err != nil {
			return nil, err
		}

		result = append(result, item)
	}
	return result, nil
}

func (dbs DBStore) SaveData(ctx context.Context, id string, value string) error {
	if id == "" || value == "" {
		return fmt.Errorf("key or value is empty")
	}
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := dbs.db.ExecContext(dbCtx, "INSERT INTO url_shortener(short_url, original_url) VALUES ($1, $2)", id, value)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			id, repeatingError := postgresInstant.getShortURLByOriginalURL(ctx, value)
			if repeatingError != nil {
				return repeatingError
			}
			return &DBError{
				ShortURL: id,
				Err:      pgErr,
			}
		} else {
			return err
		}
	}
	return nil
}

func CustomPing() bool {
	if postgresInstant.db == nil {
		return false
	}
	err := postgresInstant.db.Ping()
	return err == nil
}

func createTables(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS url_shortener (
		uuid SERIAL PRIMARY KEY,
		short_url VARCHAR NOT NULL UNIQUE,
		original_url VARCHAR NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		return err
	}
	logger.Log.Info("Database table created")
	return nil
}

func (dbs DBStore) getShortURLByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := dbs.db.QueryRowContext(dbCtx, "SELECT short_url FROM url_shortener WHERE original_url = $1", originalURL)
	var original string
	if err := row.Scan(&original); err != nil {
		return "", err
	}
	return original, nil
}
