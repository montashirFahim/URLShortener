package cfg

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// App wraps all the config needed by
// all the components to run the service
type App struct {
	ENV      ENV
	Postgres Postgres
	Redis    Redis
	Self     Self
}

// Load loads all config
func (a *App) Load() error {
	viper.AutomaticEnv()

	if err := a.ENV.Load(); err != nil {
		return err
	}

	viper.AddRemoteProvider("consul", a.ENV.ConsulURL.String(), a.ENV.ConsulPath.String())
	viper.SetConfigType("yml")
	if err := viper.ReadRemoteConfig(); err != nil {
		return err
	}
	a.Postgres.Load()
	a.Redis.Load()
	a.Self.Load()
	return nil
}

// Print prints all config
func (a *App) Print() {
	a.Postgres.Print()
	a.Redis.Print()
	a.ENV.Print()
	a.Self.Prt()
}

// NewApp returns instance of App
func NewApp() *App {
	return &App{
		Self:     Self{},
		ENV:      ENV{},
		Postgres: Postgres{},
		Redis:    Redis{},
	}
}
