package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open creates a gorm.DB connection using the supplied driver/dsn pair.
func Open(driver, dsn string) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	switch driver {
	case "sqlite", "":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	default:
		return nil, fmt.Errorf("unsupported database driver %q", driver)
	}
	if err != nil {
		return nil, err
	}

	return db, nil
}
