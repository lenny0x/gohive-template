package mysql

import (
	"database/sql"
	"fmt"

	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gohive/core/config"
)

// DSN builds the MySQL DSN used by both database/sql and GORM connections.
func DSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}

// OpenSQL opens a *sql.DB using the database/sql driver (used by goose).
func OpenSQL(cfg config.DatabaseConfig) (*sql.DB, error) {
	return sql.Open("mysql", DSN(cfg))
}

// OpenGorm opens a GORM connection.
func OpenGorm(cfg config.DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open(gormmysql.Open(DSN(cfg)), &gorm.Config{})
}

// CloseGorm best-effort closes the underlying connection of a GORM DB.
func CloseGorm(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	_ = sqlDB.Close()
}
