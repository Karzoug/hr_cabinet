package model

import (
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/stdlib" // use as driver for sqlx
	"github.com/jmoiron/sqlx"
)

const (
	defaultMaxOpenConn     = 4
	defaultMaxIdleConn     = 4
	defaultMaxConnIdleTime = time.Second * 30
	defaultMaxConnLifetime = time.Minute * 2

	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

// DB структура с настройками подключения к БД и доступом к текущему соединению.
type DB struct {
	*sqlx.DB

	maxOpenConn     int
	maxIdleConn     int
	connAttempts    int
	maxConnIdleTime time.Duration
	maxConnLifetime time.Duration
	connTimeout     time.Duration
}

// New создаёт объект DB с заданными параметрами и подключается к БД.
func New(dsn string, opts ...Option) (*DB, error) {
	db := &DB{
		maxOpenConn:     defaultMaxOpenConn,
		maxIdleConn:     defaultMaxIdleConn,
		maxConnIdleTime: defaultMaxConnIdleTime,
		maxConnLifetime: defaultMaxConnLifetime,
		connAttempts:    defaultConnAttempts,
		connTimeout:     defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(db)
	}

	var err error

	for db.connAttempts > 0 {
		if db.DB, err = sqlx.Connect("pgx", dsn); err == nil {
			break
		}

		slog.Info(fmt.Sprintf("trying to connect: attempts left %v", db.connAttempts))

		time.Sleep(db.connTimeout)

		db.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("new db: %w", err)
	}

	db.SetMaxOpenConns(db.maxOpenConn)
	db.SetMaxIdleConns(db.maxIdleConn)
	db.SetConnMaxIdleTime(db.maxConnIdleTime)
	db.SetConnMaxLifetime(db.maxConnLifetime)

	return db, nil
}

// Shutdown дожидается завершения запросов и закрывает все открытые соединения.
func (db *DB) Shutdown() error {
	return db.Close()
}
