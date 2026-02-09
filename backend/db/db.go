package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(appDataDir string) error {
	dbPath := filepath.Join(appDataDir, "sidekick.db")
	err := os.MkdirAll(appDataDir, 0755)
	if err != nil {
		return err
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Warn,
		},
	)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	// Enable WAL mode for better concurrency
	if err := db.Exec("PRAGMA journal_mode = WAL;").Error; err != nil {
		log.Printf("[db] PRAGMA journal_mode=WAL not supported or failed; continuing without WAL: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		log.Printf("[db] PRAGMA foreign_keys=ON not supported or failed; continuing without foreign key enforcement: %v", err)
	}

	DB = db
	return nil
}
