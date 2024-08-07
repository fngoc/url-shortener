package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var PostgresInstant *sql.DB

func InitializeDB(dbConf string) error {
	pqx, err := sql.Open("pgx", dbConf)
	if err != nil {
		return err
	}
	PostgresInstant = pqx
	return nil
}
