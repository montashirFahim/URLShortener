package infrastructure

import (
	"User/internal/domain/cfg"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

// Redis stors Redis client instance
type Redis struct {
	c *redis.Client
}

// Connect creates a Redis client connection
func (r *Redis) Connect(cnf *cfg.App) error {

	rc := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cnf.Redis.Host.String(), cnf.Redis.Port.Int()),
		DB:   cnf.Redis.Db.Int(), // use default DB
	})

	_, err := rc.Ping().Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Redis Client Connection")
		return err
	}

	r.c = rc

	return nil
}

// Client returns unexported *redis.Client
func (r *Redis) Client() *redis.Client {
	return r.c
}
