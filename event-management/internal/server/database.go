package server

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDB(cfg *Postgres) *gorm.DB {
	dsn := cfg.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	pq, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database connection: %v", err)
	}

	pq.SetMaxOpenConns(cfg.MaxOpenConns)
	pq.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	pq.SetMaxIdleConns(cfg.MaxIdleConns)
	pq.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db
}
