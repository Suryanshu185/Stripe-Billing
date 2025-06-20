package db

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	dbModels "gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	logger := viper.Get("Logger").(*zap.Logger)
	logger.Debug("connecting to database")
	dbHost, dbPort, dbName, dbUserName, dbPassword := viper.GetString("DB_HOST"), viper.GetString("DB_PORT"), viper.GetString("DB_NAME"),
		viper.GetString("DB_POSTGRES_USERNAME"), viper.GetString("DB_POSTGRES_PASSWORD")
	err := createDatabaseIfNotExists(dbHost, dbPort, dbName, dbUserName, dbPassword)
	if err != nil {
		log.Fatalf("error creating db: %v", err)
	}

	dsn := fmt.Sprintf("host=%s dbname=%s port=%s", dbHost, dbName, dbPort)
	if len(dbUserName) > 0 && len(dbPassword) > 0 {
		dsn = dsn + fmt.Sprintf(" user=%s password=%s", dbUserName, dbPassword)
	}
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("error connecting to database", zap.Error(err))
		log.Fatalf("Failed to connect to database: %v", err)
	}
	runMigrations()
}

func runMigrations() {
	err := DB.AutoMigrate(
		&dbModels.Account{},
		&dbModels.ActiveResource{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func createDatabaseIfNotExists(dbHost, dbPort, dbName, dbUserName, dbPassword string) error {
	logger := viper.Get("Logger").(*zap.Logger)
	dsn := fmt.Sprintf("host=%s dbname=%s port=%s", dbHost, dbName, dbPort)
	if len(dbUserName) > 0 && len(dbPassword) > 0 {
		dsn = dsn + fmt.Sprintf(" user=%s password=%s", os.Getenv("DB_POSTGRES_USERNAME"), os.Getenv("DB_POSTGRES_PASSWORD"))
	}
	logger.Info("Connecting to database", zap.String("database host", dbHost), zap.String("database name", dbName))
	tempDbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("error connecting to database", zap.Error(err))
		return err
	}

	sqlDB, err := tempDbConn.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	// Check if the database "infra" exists
	rows, err := sqlDB.Query(fmt.Sprintf("SELECT datname FROM pg_database WHERE datname = '%s'", dbName))
	if err != nil {
		logger.Error("error querying database", zap.Error(err))
		return err
	}
	defer rows.Close()

	// If the database does not exist, create it
	if !rows.Next() {
		logger.Info("Database does not exist. Creating it.")
		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
