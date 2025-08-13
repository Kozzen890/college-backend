package db

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	instance     *gorm.DB
	instanceOnce sync.Once
)

// Connect initializes the global DB connection (Postgres if DSN looks like postgres URL, otherwise SQLite path)
func Connect(dsnOrSqlitePath string) error {
	var err error
	instanceOnce.Do(func() {
		var openErr error
		lower := strings.ToLower(dsnOrSqlitePath)

		// Configure GORM with logging
		config := &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}

		log.Printf("Attempting to connect to database: %s", dsnOrSqlitePath)

		switch {
		case strings.HasPrefix(lower, "postgres://") || strings.HasPrefix(lower, "postgresql://"):
			log.Printf("Using PostgreSQL driver")
			instance, openErr = gorm.Open(postgres.Open(dsnOrSqlitePath), config)
		case strings.HasPrefix(lower, "mysql://") || strings.Contains(lower, "@tcp("):
			log.Printf("Using MySQL driver")
			instance, openErr = gorm.Open(mysql.Open(dsnOrSqlitePath), config)
		default:
			log.Printf("Using SQLite driver")
			instance, openErr = gorm.Open(sqlite.Open(dsnOrSqlitePath), config)
		}

		if openErr != nil {
			err = fmt.Errorf("open database: %w", openErr)
			log.Printf("Database connection failed: %v", err)
			return
		}

		log.Printf("Database connection successful")

		// Test the connection
		sqlDB, sqlErr := instance.DB()
		if sqlErr != nil {
			err = fmt.Errorf("get sql.DB: %w", sqlErr)
			log.Printf("Failed to get sql.DB: %v", err)
			return
		}

		if pingErr := sqlDB.Ping(); pingErr != nil {
			err = fmt.Errorf("ping database: %w", pingErr)
			log.Printf("Database ping failed: %v", err)
			return
		}

		log.Printf("Database ping successful")
	})
	return err
}

// Instance returns the initialized DB instance
func Instance() *gorm.DB {
	return instance
}

// ConnectMySQL creates a direct MySQL connection for testing
func ConnectMySQL(host, port, dbname, username, password string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	log.Printf("MySQL DSN: %s", dsn)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	log.Printf("MySQL connection successful")
	return db, nil
}
