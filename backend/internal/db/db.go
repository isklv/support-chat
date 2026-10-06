package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init connects to database (SQLite or PostgreSQL) and runs auto migrations.
func Init(driver, dsn string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	driver = strings.ToLower(strings.TrimSpace(driver))
	if driver == "postgres" || strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		dialector = postgres.Open(dsn)
	} else {
		// Default to SQLite
		if dsn == "" {
			dsn = "./data/support.db"
		}
		dir := filepath.Dir(dsn)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("create db directory %s: %w", dir, err)
			}
		}
		dialector = sqlite.Open(dsn)
	}

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	database, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// Auto migrate tables
	if err := database.AutoMigrate(
		&SupportThread{},
		&SupportMessage{},
		&SupportAttachment{},
		&SupportToken{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	log.Printf("db: connected and migrated successfully")
	return database, nil
}
