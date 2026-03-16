package conn

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"magic.pathao.com/parcel/prism/internal/config"
)

func ConnectPostgres(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	pg := cfg.Postgres
	dsn := pg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(pg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(pg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(pg.ConnMaxLifetime) * time.Second)

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("successfully connected to database")
	return db, nil
}
