package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"easy-arbitra/backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Migration struct {
	Filename  string `gorm:"primaryKey;size:255"`
	AppliedAt time.Time
}

func (Migration) TableName() string { return "schema_migration" }

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}

	if err := db.AutoMigrate(&Migration{}); err != nil {
		log.Fatalf("ensure migration table: %v", err)
	}

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatalf("scan migrations: %v", err)
	}
	sort.Strings(files)
	if len(files) == 0 {
		log.Println("no migration files found")
		return
	}

	ctx := context.Background()
	for _, file := range files {
		name := filepath.Base(file)
		applied, err := isApplied(ctx, db, name)
		if err != nil {
			log.Fatalf("check migration %s: %v", name, err)
		}
		if applied {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("read migration %s: %v", name, err)
		}
		sql := strings.TrimSpace(string(content))
		if sql == "" {
			if err := markApplied(ctx, db, name); err != nil {
				log.Fatalf("mark empty migration %s: %v", name, err)
			}
			continue
		}

		if err := db.WithContext(ctx).Exec(sql).Error; err != nil {
			log.Fatalf("apply migration %s: %v", name, err)
		}
		if err := markApplied(ctx, db, name); err != nil {
			log.Fatalf("mark migration %s: %v", name, err)
		}
		fmt.Printf("applied %s\n", name)
	}

	log.Println("migrations complete")
}

func isApplied(ctx context.Context, db *gorm.DB, filename string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(&Migration{}).Where("filename = ?", filename).Count(&count).Error
	return count > 0, err
}

func markApplied(ctx context.Context, db *gorm.DB, filename string) error {
	return db.WithContext(ctx).Create(&Migration{Filename: filename, AppliedAt: time.Now().UTC()}).Error
}
