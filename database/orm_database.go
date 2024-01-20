package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/deveusss/evergram-core/caching"
	"github.com/deveusss/evergram-core/config"

	"github.com/eapache/go-resiliency/retrier"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// OrmDatabase represents a wrapper around gorm.DB
type OrmDatabase struct {
	Orm           *gorm.DB
	Cache         caching.AppCacher
	Retry         *retrier.Retrier
	EnableCaching bool // Flag to enable or disable caching for queries
	ctx           context.Context
	slog          *slog.Logger
}

func New(log *slog.Logger, config *config.DatabaseConfig) (*OrmDatabase, error) {
	return NewWithContext(nil, config, false, log)
}

// NewDatabaseWithContext creates a new instance of OrmDatabase with context
func NewWithContext(ctx context.Context, config *config.DatabaseConfig, enableCaching bool, slog *slog.Logger) (*OrmDatabase, error) {
	dsn := buildConnectionString(config)
	slog.Info("Connecting to database", "dsn", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Error opening database: %v", err)
		return nil, err
	}
	inner, err := db.DB()

	if ctx != nil {
		// Check the connection during creation
		if inner.PingContext(ctx); err != nil {
			slog.Error("Error pinging database: %v", err)
			return nil, err
		}
	}

	// Configure connection pool
	inner.SetMaxOpenConns(config.MaxOpenConns)
	inner.SetMaxIdleConns(config.MaxIdleConns)

	// Initialize cache
	cache, err := caching.NewAppCache()
	if err != nil {
		slog.Error("Error initializing cache: %v", err)
		return nil, err
	}

	retry := retrier.New(retrier.ExponentialBackoff(config.MaxRetries, config.RetryWait), nil)
	slog.Info("Connected to", "dsn", dsn)
	return &OrmDatabase{ctx: ctx, Orm: db, Cache: cache, Retry: retry, EnableCaching: enableCaching}, nil
}

// OpenConnection opens a connection to the database
func (db *OrmDatabase) OpenConnection() error {
	innerDb, err := db.Orm.DB()
	if err != nil {
		db.slog.Error("Failed when getting inner DB: %v", err)
		return err
	}
	return innerDb.Ping()
}

// CloseConnection closes the connection to the database
func (db *OrmDatabase) CloseConnection() error {
	innerDb, err := db.Orm.DB()
	if err != nil {
		db.slog.Error("Failed when getting inner DB: %v", err)
		return err
	}
	return innerDb.Close()
}

// WithTransactionWithContext executes a function inside a transaction with context
func (db *OrmDatabase) WithTransactionContext(context context.Context, fn func(*OrmDatabase) error) error {
	tx := db.Orm.WithContext(context).Begin()
	if tx.Error != nil {
		db.slog.Error("Error beginning transaction: %v", tx.Error)
		return tx.Error
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			db.slog.Error("Error in transaction: %v", p)
		} else if err := recover(); err != nil {
			_ = tx.Rollback()
		} else {
			err := tx.Commit().Error
			if err != nil {
				_ = tx.Rollback()
				db.slog.Error("Error committing transaction: %v", err)
			}
		}
	}()

	return fn(&OrmDatabase{tx, db.Cache, db.Retry, db.EnableCaching, context, db.slog})
}

// buildConnectionString builds the database connection string for PostgreSQL
func buildConnectionString(dbConfig *config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)
}

func (db *OrmDatabase) AuthMigrate(dst ...interface{}) error {
	return db.Orm.AutoMigrate(dst...)
}
