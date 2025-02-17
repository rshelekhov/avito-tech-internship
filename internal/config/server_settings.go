package config

import "github.com/rshelekhov/merch-store/internal/config/settings"

type ServerSettings struct {
	AppEnv       string                `mapstructure:"APP_ENV"`
	HTTPServer   settings.HTTPServer   `mapstructure:",squash"`
	Postgres     settings.Postgres     `mapstructure:",squash"`
	JWT          settings.JWT          `mapstructure:",squash"`
	PasswordHash settings.PasswordHash `mapstructure:",squash"`
}
