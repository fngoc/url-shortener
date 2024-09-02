package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/fngoc/url-shortener/cmd/shortener/config"
	"github.com/fngoc/url-shortener/cmd/shortener/constants"
	"github.com/fngoc/url-shortener/internal/logger"
	"github.com/fngoc/url-shortener/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	"sync"
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

type DBDeleteError struct {
	Message string
}

func (d *DBDeleteError) Error() string {
	return d.Message
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

	row := dbs.db.QueryRowContext(dbCtx, "SELECT original_url, is_deleted FROM url_shortener WHERE short_url = $1", key)
	var originalURL string
	var deleteFlag bool

	err := row.Scan(&originalURL, &deleteFlag)
	if err != nil {
		return "", err
	}

	if deleteFlag {
		return "", &DBDeleteError{
			Message: "shortener is already deleted",
		}
	}

	return originalURL, nil
}

func (dbs DBStore) GetAllData(ctx context.Context) ([]models.ResponseDto, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	userID := ctx.Value(constants.UserIDKey).(int)
	rows, err := dbs.db.QueryContext(dbCtx, "SELECT short_url, original_url FROM url_shortener WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var result []models.ResponseDto
	for rows.Next() {
		var shortURL string
		var originalURL string

		if err := rows.Scan(&shortURL, &originalURL); err != nil {
			return nil, err
		}

		result = append(result, models.ResponseDto{
			ShortURL:    config.Flags.BaseResultAddress + "/" + shortURL,
			OriginalURL: originalURL,
		})
	}
	return result, nil
}

func (dbs DBStore) SaveData(ctx context.Context, id string, value string) error {
	if id == "" || value == "" {
		return fmt.Errorf("key or value is empty")
	}
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	userID := ctx.Value(constants.UserIDKey).(int)
	_, err := dbs.db.ExecContext(dbCtx, "INSERT INTO url_shortener(short_url, original_url, user_id) VALUES ($1, $2, $3)", id, value, userID)
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

func (dbs DBStore) DeleteData(rCtx context.Context, userID int, urls []string) error {
	if len(urls) == 0 {
		return nil
	}

	batchSize := 10
	errChan := make(chan error, 1)
	ctx, cancel := context.WithTimeout(rCtx, 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	for i := 0; i < len(urls); i += batchSize {
		end := i + batchSize
		if end > len(urls) {
			end = len(urls)
		}
		batchIDs := urls[i:end]

		wg.Add(1)
		go func(batchIDs []string) {
			defer wg.Done()

			query := "UPDATE url_shortener SET is_deleted = true WHERE user_id = $1 AND short_url = ANY($2::text[])"
			_, err := dbs.db.ExecContext(ctx, query, userID, pq.Array(batchIDs))
			if err != nil {
				select {
				case errChan <- fmt.Errorf("delete batch error: %v", err):
					cancel()
				default:
				}
			}
		}(batchIDs)
	}

	go func() {
		wg.Wait()
		close(errChan)

		select {
		case err := <-errChan:
			logger.Log.Warn(err.Error())
		default:
		}
	}()

	return nil
}

func createTables(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS url_shortener (
		uuid SERIAL PRIMARY KEY,
		short_url VARCHAR NOT NULL UNIQUE,
		original_url VARCHAR NOT NULL UNIQUE,
		user_id BIGSERIAL NOT NULL,
		is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	createIndexQuery := `CREATE UNIQUE INDEX IF NOT EXISTS short_url_idx ON url_shortener (short_url)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		return err
	}
	_, errIdx := db.ExecContext(ctx, createIndexQuery)
	if errIdx != nil {
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
