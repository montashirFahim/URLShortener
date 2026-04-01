package cfg

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Postgres struct {
	host            StrVal
	port            StrVal
	user            StrVal
	password        StrVal
	dbName          StrVal
	MaxIdleConn     IntVal
	MaxOpenConn     IntVal
	MaxConnLifetime IntVal
}

// Load loads all the Postgres configuration
func (pg *Postgres) Load() {
	pg.host.Load("postgres.host")
	pg.port.Load("postgres.port")
	pg.user.Load("postgres.user")
	pg.password.Load("postgres.password")
	pg.dbName.Load("postgres.db_name")

	pg.MaxIdleConn.Load("postgres.maxidleconn")
	pg.MaxOpenConn.Load("postgres.maxopenconn")
	pg.MaxConnLifetime.Load("postgres.maxconnlifetime")
}

func (pg *Postgres) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pg.host,
		pg.port,
		pg.user,
		pg.password,
		pg.dbName,
	)
}

func (pg *Postgres) Print() {
	logrus.Info("Postgres-host: ", pg.host)
	logrus.Info("Postgres-port: ", pg.port)
	logrus.Info("Postgres-user: ", pg.user)
	logrus.Info("Postgres-password: ", pg.password)
	logrus.Info("Postgres-db: ", pg.dbName)

	logrus.Info("Postgres-MaxIdleConn: ", pg.MaxIdleConn)
	logrus.Info("Postgres-MaxOpenConn: ", pg.MaxOpenConn)
	logrus.Info("Postgres-MaxConnLifetime: ", pg.MaxConnLifetime)
}
