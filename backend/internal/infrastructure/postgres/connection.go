package postgres

import (
	"fmt"
	"log"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.ConnectionString()
	
	log.Printf("Connecting to database: host=%s port=%s user=%s dbname=%s", 
		cfg.Host, cfg.Port, cfg.User, cfg.DBName)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database at %s:%s (user: %s, db: %s): %w. "+
			"Make sure PostgreSQL is running and accessible", 
			cfg.Host, cfg.Port, cfg.User, cfg.DBName, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database at %s:%s: %w. "+
			"Check if PostgreSQL is running and the connection settings are correct", 
			cfg.Host, cfg.Port, err)
	}

	log.Printf("Successfully connected to database: %s@%s:%s/%s", 
		cfg.User, cfg.Host, cfg.Port, cfg.DBName)
	return db, nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

