package config

import (
	"path/filepath"
	"strings"

	l "github.com/deesel/wol/internal/logger"
	"github.com/spf13/viper"
)

// ServerConfig hold API server configuration
type ServerConfig struct {
	Address string
	Port    int
}

// APIKeyConfig hold API key configuration
type APIKeyConfig struct {
	Name string
	Key  string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Enabled bool
	APIKeys []APIKeyConfig
}

// Config holds service configuration
type Config struct {
	Server ServerConfig
	Auth   AuthConfig
}

var defaults map[string]interface{} = map[string]interface{}{
	"server.address": "0.0.0.0",
	"server.port":    8001,
	"auth.enabled":   false,
}

// New creates new configuration instance composed of default values overridden by values specified in configuration file
func New(file string) (*Config, error) {
	c := &Config{}
	v := viper.New()
	configname := strings.Trim(filepath.Base(file), filepath.Ext(file))
	dirname := filepath.Dir(file)

	v.SetConfigName(configname)
	v.AddConfigPath(dirname)
	v.AddConfigPath("/etc/wol")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	for key, val := range defaults {
		v.SetDefault(key, val)
	}

	err := v.ReadInConfig()
	if err != nil {
		l.New().Logger().Warn(err)
	}

	err = v.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
