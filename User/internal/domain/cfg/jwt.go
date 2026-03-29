package cfg

import "github.com/sirupsen/logrus"

type JWT struct {
	Secret StrVal
}

func (j *JWT) Load() {
	j.Secret.Load("jwt.secret")
}

func (j *JWT) Print() {
	logrus.Info("JWT-secret : ", j.Secret)
}
