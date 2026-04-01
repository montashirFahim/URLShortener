package cfg

import (
	"github.com/sirupsen/logrus"
)

// Redis holds configuration of Redis
type Redis struct {
	Host StrVal
	Port IntVal
	Db   IntVal
}

// Load loads all the redis configuration
func (r *Redis) Load() {
	r.Host.Load("redis.host")
	r.Port.Load("redis.port")
	r.Db.Load("redis.db")
}

// Print prints all the redis configuration
func (r *Redis) Print() {
	logrus.Info("Redis-host : ", r.Host)
	logrus.Info("Redis-port : ", r.Port)
	logrus.Info("Redis-db : ", r.Db)
}
