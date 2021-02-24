package storage

import (
	"fmt"
	"github.com/fatihsezgin/candlecloud-backend/internal/config"
	"github.com/fatihsezgin/candlecloud-backend/internal/storage/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	db    *gorm.DB
	users UserRepository
}

func DBConn(cfg *config.DatabaseConfiguration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	db, err = gorm.Open("postgres", "host="+cfg.Host+" port="+cfg.Port+" user="+cfg.Username+" dbname="+cfg.Name+"  sslmode=disable password="+cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("could not open postgresql connection: %w", err)
	}

	db.LogMode(cfg.LogMode)

	return db, err
}

func New(db *gorm.DB) *Database {
	return &Database{
		db:    db,
		users: user.NewRepository(db),
	}
}

// Users returns the UserRepository.
func (db *Database) Users() UserRepository {
	return db.users
}
