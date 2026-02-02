package postgres

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func RunMigrations(db *gorm.DB, migrationsDir string) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	if len(sqlFiles) == 0 {
		return nil
	}

	sort.Strings(sqlFiles)

	silentDB := db.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	for _, fileName := range sqlFiles {
		filePath := filepath.Join(migrationsDir, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		if err := silentDB.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}
	}

	fmt.Printf("Migrations successfully applied (%d files)\n", len(sqlFiles))
	return nil
}

