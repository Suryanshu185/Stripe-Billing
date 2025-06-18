package repository

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/viper"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/db"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//TODO : do we want to wrap the original error ?

func GetDBError(err error) error {
	logger := viper.Get("Logger").(*zap.Logger)
	logger.Info("db error", zap.Error(err))
	var pgError *pgconn.PgError
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return db.ErrRecordNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return db.ErrDuplicatedRecord
	case errors.As(err, &pgError):
		if pgError.Code == "23505" {
			return db.ErrDuplicatedRecord
		} else {
			return db.ErrInternalServer
		}
	default:
		return db.ErrInternalServer
	}
}
