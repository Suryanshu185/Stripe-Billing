package db

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	dbModels "gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	logger := viper.Get("Logger").(*zap.Logger)
	logger.Debug("Initializing database...")

	// Railway sets this automatically
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Manual fallback (optional for local dev)
		dbHost := os.Getenv("PGHOST")
		dbPort := os.Getenv("PGPORT")
		dbName := os.Getenv("PGDATABASE")
		dbUser := os.Getenv("PGUSER")
		dbPassword := os.Getenv("PGPASSWORD")

		if dbHost == "" || dbPort == "" || dbName == "" || dbUser == "" || dbPassword == "" {
			log.Fatal("Database connection details missing in environment variables")
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
			dbHost, dbUser, dbPassword, dbName, dbPort)
	}

	logger.Info("Connecting to database...", zap.String("dsn", dsn))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		log.Fatalf("Failed to connect to database: %v", err)
	}

	logger.Info("Connected to DB. Running migrations...")
	runMigrations()
	logger.Info("Migrations done.")
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
