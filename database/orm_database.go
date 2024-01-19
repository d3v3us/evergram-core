package database

import (
	"context"
	"fmt"
	"log"

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
}

func New(config *config.DatabaseConfig) (*OrmDatabase, error) {
	return NewWithContext(context.Background(), config, false)
}

// NewDatabaseWithContext creates a new instance of OrmDatabase with context
func NewWithContext(ctx context.Context, config *config.DatabaseConfig, enableCaching bool) (*OrmDatabase, error) {
	dsn := buildConnectionString(config)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}
	inner, err := db.DB()

	// Check the connection during creation
	if inner.PingContext(ctx); err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	// Configure connection pool
	inner.SetMaxOpenConns(10)
	inner.SetMaxIdleConns(5)

	// Initialize cache
	cache, err := caching.NewAppCache()
	if err != nil {
		log.Printf("Error initializing cache: %v", err)
		return nil, err
	}

	retry := retrier.New(retrier.ExponentialBackoff(config.MaxRetries, config.RetryWait), nil)

	return &OrmDatabase{ctx: ctx, Orm: db, Cache: cache, Retry: retry, EnableCaching: enableCaching}, nil
}

// OpenConnection opens a connection to the database
func (db *OrmDatabase) OpenConnection() error {
	innerDb, err := db.Orm.DB()
	if err != nil {
		log.Printf("Failed when getting inner DB: %v", err)
		return err
	}
	return innerDb.Ping()
}

// CloseConnection closes the connection to the database
func (db *OrmDatabase) CloseConnection() error {
	innerDb, err := db.Orm.DB()
	if err != nil {
		log.Printf("Failed when getting inner DB: %v", err)
		return err
	}
	return innerDb.Close()
}

// WithTransactionWithContext executes a function inside a transaction with context
func (db *OrmDatabase) WithTransactionContext(context context.Context, fn func(*OrmDatabase) error) error {
	tx := db.Orm.WithContext(context).Begin()
	if tx.Error != nil {
		log.Printf("Error beginning transaction: %v", tx.Error)
		return tx.Error
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			log.Panic(p)
		} else if err := recover(); err != nil {
			_ = tx.Rollback()
		} else {
			err := tx.Commit().Error
			if err != nil {
				_ = tx.Rollback()
				log.Printf("Error committing transaction: %v", err)
			}
		}
	}()

	return fn(&OrmDatabase{tx, db.Cache, db.Retry, db.EnableCaching, context})
}

// buildConnectionString builds the database connection string for PostgreSQL
func buildConnectionString(dbConfig *config.DatabaseConfig) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
}
