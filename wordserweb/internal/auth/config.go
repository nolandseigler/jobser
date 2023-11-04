package auth

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	PubKeyPath  string `mapstructure:"JWT_PUB_KEY_PATH"`
	PrivKeyPath string `mapstructure:"JWT_PRIV_KEY_PATH"`
}

func ConfigFromEnv() (Config, error) {
	c := Config{}
	if err := viper.BindEnv("JWT_PUB_KEY_PATH"); err != nil {
		return c, fmt.Errorf("failed to bind 'JWT_PUB_KEY_PATH'")
	}
	viper.SetDefault("JWT_PUB_KEY_PATH", "/etc/wordserweb/keys/pub.rsa.pem")

	if err := viper.BindEnv("JWT_PRIV_KEY_PATH"); err != nil {
		return c, fmt.Errorf("failed to bind 'JWT_PRIV_KEY_PATH'")
	}
	viper.SetDefault("JWT_PRIV_KEY_PATH", "/etc/wordserweb/keys/priv.rsa.pem")

	if err := viper.Unmarshal(&c); err != nil {
		return c, fmt.Errorf("failed to unmarshal config")
	}

	return c, nil
}
