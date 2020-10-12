package config

import (
	. "github.com/MrNinso/MyGoToolBox/lang/env"
)

type sslConfig struct {
	CertPath string
	KeyPath  string
}

type Config struct {
	SSL        *sslConfig
	Bind       string
	BcryptCost int
	LogPath    string
}

var Configs = LoadConfigs()

func LoadConfigs() *Config {
	return &Config{
		Bind:       GetEnv("BIND", "0.0.0.0:1443"),
		BcryptCost: GetIntFromEnv("BCOST", 12),
		SSL: &sslConfig{
			CertPath: GetEnv("CERTPATH", "./certs/server.crt"),
			KeyPath:  GetEnv("KEYPATH", "./certs/server.key"),
		},
		LogPath: GetEnv("LOGPATH", "/var/log/Lab"),
	}
}
