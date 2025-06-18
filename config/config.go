package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

func LoadConfig() {
	//read configuration parameters from environment
	viper.AutomaticEnv()

	setDefaults()
	err := initLogger()
	if err != nil {
		panic(fmt.Errorf("fatal configuring logger: %w", err))
	}

	err = checkRequiredKeys()
	if err != nil {
		panic(fmt.Errorf("invalid config file : %w", err))
	}
}

func setDefaults() {
	if viper.GetString("OWN_PORT") == "" {
		viper.Set("OWN_PORT", "8800")
	}
	if viper.GetString("DB_HOST") == "" {
		viper.Set("DB_HOST", "localhost")
	}
	if viper.GetString("DB_PORT") == "" {
		viper.Set("DB_PORT", "5432")
	}
	if viper.GetString("DB_NAME") == "" {
		viper.Set("DB_NAME", "vcs")
	}
	if viper.GetString("LOG_LEVEL") == "" {
		viper.Set("LOG_LEVEL", "info")
	}

	if viper.GetString("ENVIRONMENT") == "" {
		viper.Set("ENVIRONMENT", "development")
	}
}

func initLogger() error {
	var logLevel zapcore.Level

	switch viper.GetString("LOG_LEVEL") {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	case "fatal":
		logLevel = zapcore.FatalLevel
	default:
		logLevel = zapcore.InfoLevel // Default to info level if not set
	}

	// Create a logger configuration
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Build the logger
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v", err)
		return err
	}

	// Assign the logger to the config
	viper.Set("Logger", logger)
	return nil
}

func checkRequiredKeys() error {
	requiredKeys := []string{
		"OWN_PORT",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"LOG_LEVEL",
	}
	for _, key := range requiredKeys {
		if !viper.IsSet(key) {
			return fmt.Errorf("missing required configuration key: %s", key)
		}
	}
	return nil
}
