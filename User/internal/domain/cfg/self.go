package cfg

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Self stores config specific to the Self
type Self struct {
	port  int
	debug bool
	mu    sync.RWMutex
}

// Prt prints values
func (s *Self) Prt() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	logrus.WithFields(logrus.Fields{
		"port":  s.port,
		"debug": s.debug,
	}).Info("self")
}

// Port returns port
func (s *Self) Port() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.port
}

// Debug retruns debug flag
func (s *Self) Debug() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.debug
}

// Load loads app configuration
func (s *Self) Load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.debug = viper.GetBool("debug")
	s.port = viper.GetInt("port")
}
