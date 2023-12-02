package auth

import (
	"errors"

	serverConfig "github.com/AlyonaAg/bot-detector/internal/config"
)

var (
	noKeyEnvironmentVariables = errors.New("no key in environment variables")
)

type Config struct {
	signingKey     string
	expireDuration int64
}

func NewConfig() (*Config, error) {
	expireDuration, err := serverConfig.GetValue(serverConfig.ExpireDuration)
	if err != nil {
		return nil, err
	}

	signingKey, err := serverConfig.GetValue(serverConfig.SigningKey)
	if err != nil {
		return nil, err
	}

	return &Config{
		signingKey:     signingKey.(string),
		expireDuration: expireDuration.(int64),
	}, nil
}
