package db

import (
	"fmt"
	"localVercel/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "novexdeploy"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		host, port, user, password, dbname)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Автоматическая миграция
	err = DB.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.EnvVar{},
		&models.Deployment{},
		&models.Session{},
	)
	
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("✅ Database connected and migrated")
	return nil
}