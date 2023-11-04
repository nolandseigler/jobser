package postgres

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Protocol     string `mapstructure:"POSTGRES_PROTOCOL"`
	Username     string `mapstructure:"POSTGRES_USERNAME"`
	Password     string `mapstructure:"POSTGRES_PASSWORD"`
	Hostname     string `mapstructure:"POSTGRES_HOSTNAME"`
	Port         int    `mapstructure:"POSTGRES_PORT"`
	DatabaseName string `mapstructure:"POSTGRES_DATABASE_NAME"`
}

func ConfigFromEnv() (Config, error) {
	c := Config{}
	if err := viper.BindEnv("POSTGRES_PROTOCOL"); err != nil {
		return c, fmt.Errorf("failed to bind 'POSTGRES_PROTOCOL'")
	}
	viper.SetDefault("POSTGRES_PROTOCOL", "postgres")

	if err := viper.BindEnv("POSTGRES_USERNAME"); err != nil {
		return c, fmt.Errorf("failed to bind 'POSTGRES_USERNAME'")
	}
	viper.SetDefault("POSTGRES_USERNAME", "wordser")

	if err := viper.BindEnv("POSTGRES_PASSWORD"); err != nil {
		return c, fmt.Errorf("failed to bind 'POSTGRES_PASSWORD'")
	}

	if err := viper.BindEnv("POSTGRES_HOSTNAME"); err != nil {
		return c, fmt.Errorf("failed to bind 'POSTGRES_HOSTNAME'")
	}
	viper.SetDefault("POSTGRES_HOSTNAME", "db")

	if err := viper.BindEnv("POSTGRES_PORT"); err != nil {
		return c, fmt.Errorf("failed to bind 'POSTGRES_PORT'")
	}
	viper.SetDefault("POSTGRES_PORT", 5432)

	if err := viper.BindEnv("POSTGRES_DATABASE_NAME"); err != nil {
		return c, fmt.Errorf("failed to bind 'POSTGRES_DATABASE_NAME'")
	}
	viper.SetDefault("POSTGRES_DATABASE_NAME", "wordser")

	if err := viper.Unmarshal(&c); err != nil {
		return c, fmt.Errorf("failed to unmarshal config")
	}

	return c, nil
}
