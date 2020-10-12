package config

import (
	"github.com/MrNinso/MyGoToolBox/lang/env"
)

type sslConfig struct {
	CertPath string
	KeyPath  string
}

type config struct {
	SSL        *sslConfig
	Bind       string
	BcryptCost int
	LogPath    string
}

//ConfigManager entrypoint
var Configs = loadConfigs()

func loadConfigs() *config {
	return &config{
		Bind:       env.GetEnv("BIND", "0.0.0.0:1443"),
		BcryptCost: env.GetIntFromEnv("BCOST", 12),
		SSL: &sslConfig{
			CertPath: env.GetEnv("CERTPATH", "./certs/server.crt"),
			KeyPath:  env.GetEnv("KEYPATH", "./certs/server.key"),
		},
		LogPath: env.GetEnv("LOGPATH", "/var/log/Lab"),
	}
}
