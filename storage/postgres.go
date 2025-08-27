package storage

import (
	"database/sql"
	"fmt"
	"github.com/ZaharBorisenko/jwt-auth/models"
	_ "github.com/lib/pq"
)

func ConnectDB(config *models.Config) (*sql.DB, error) {
	connectString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.HOST, config.PORT, config.USER, config.PASSWORD, config.DBNAME,
	)

	DB, err := sql.Open("postgres", connectString)
	if err != nil {
		return nil, fmt.Errorf("no connect DB %w", err)
	}
	err = DB.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot ping DB: %w", err)
	}

	return DB, nil
}
