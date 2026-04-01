package infrastructure

import (
	"Server/internal/entities/cfg"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Postgres stores sqlx.DB instance
type Postgres struct {
	DB *sqlx.DB
}

// Connect creates a postgres db connection
func (pg *Postgres) Connect(cnf *cfg.App) error {
	uri := cnf.Postgres.DSN()
	logrus.Info("postgres uri: " + uri)

	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		return err
	}

	// Set connection pool config
	db.SetMaxIdleConns(cnf.Postgres.MaxIdleConn.Int())
	db.SetMaxOpenConns(cnf.Postgres.MaxOpenConn.Int())
	db.SetConnMaxLifetime(time.Second * time.Duration(cnf.Postgres.MaxConnLifetime.Int()))

	// Test connection
	if err := db.Ping(); err != nil {
		return err
	}

	pg.DB = db
	return nil
}

// NewPostgres returns initialized Postgres instance
func NewPostgres(db *sqlx.DB) Postgres {
	return Postgres{db}
}

// Close closes the DB connection
func (pg *Postgres) Close() error {
	return pg.DB.Close()
}
