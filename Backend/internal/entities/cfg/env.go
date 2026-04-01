package cfg

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// ENV holds all env value
type ENV struct {
	ConsulURL  StrVal
	ConsulPath StrVal
}

// Load loads env values
func (e *ENV) Load() error {
	e.ConsulURL.Load("consul_url")
	if e.ConsulURL == "" {
		return errors.New("consul url is empty")
	}

	e.ConsulPath.Load("consul_path")
	if e.ConsulPath == "" {
		return errors.New("consul path is empty")
	}

	return nil
}

// Print prints all env config
func (e *ENV) Print() {
	logrus.Info("consul-url: ", e.ConsulURL)
	logrus.Info("consul-path: ", e.ConsulPath)
}
