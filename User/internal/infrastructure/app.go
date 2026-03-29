package infrastructure

import "User/internal/domain/cfg"

// App encapsulates db conns
type App struct {
	postgresql *Postgres
	redis      *Redis
}

// MySQL returns mysql db conn
func (a *App) PostgreSQL() *Postgres {
	return a.postgresql
}

// Redis returns redis client
func (a *App) Redis() *Redis {
	return a.redis
}

// Connect takes config.App and connects mongo and postgres db
func (a *App) Connect(cnf *cfg.App) error {

	if err := a.postgresql.Connect(cnf); err != nil {
		return err
	}

	if err := a.redis.Connect(cnf); err != nil {
		return err
	}

	return nil
}

// NewApp returns a connection app
func NewApp() *App {
	return &App{
		postgresql: &Postgres{},
		redis:      &Redis{},
	}
}
