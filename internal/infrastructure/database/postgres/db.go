package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
	"todo-list/config"
	"todo-list/internal/domain/model"
)

type Database interface {
	GetDB() *gorm.DB
	Close() error
}

type PostgresDB struct {
	db *gorm.DB
}

func createConnection(cfg *config.DatabaseConfig) (*PostgresDB, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	// automigrate
	if err := db.AutoMigrate(&model.Task{}, &model.Tag{}); err != nil {
		return nil, err
	}
	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) GetDB() *gorm.DB {
	return p.db
}

func (p *PostgresDB) Close() error {
	conn, err := p.db.DB()
	if err != nil {
		return err
	}
	return conn.Close()
}

var (
	instance Database
	once     sync.Once
)

func ProvideDBClient(cfg *config.DatabaseConfig) (Database, error) {
	var err error
	once.Do(func() {
		instance, err = createConnection(cfg)
	})
	return instance, err
}

func validateConfig(cfg *config.DatabaseConfig) error {
	if cfg.Host == "" || cfg.Port == 0 || cfg.User == "" || cfg.Password == "" || cfg.DBName == "" {
		return fmt.Errorf("invalid database config")
	}
	return nil
}
