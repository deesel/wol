package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address string
	Port    int
}

type APIKeyConfig struct {
	Name string
	Key  string
}

type AuthConfig struct {
	Enabled bool
	APIKeys []APIKeyConfig
}

type Config struct {
	Server ServerConfig
	Auth   AuthConfig
}

var defaults map[string]interface{} = map[string]interface{}{
	"server.address": "0.0.0.0",
	"server.port":    8001,
	"auth.enabled":   false,
}

func New(file string) (*Config, error) {
	c := &Config{}
	v := viper.New()
	configname := strings.Trim(filepath.Base(file), filepath.Ext(file))
	dirname := filepath.Dir(file)

	v.SetConfigName(configname)
	v.AddConfigPath(dirname)
	v.AddConfigPath("/etc/wol")
	v.AddConfigPath("$HOME/.wol")

	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = v.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
