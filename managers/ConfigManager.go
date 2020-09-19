package managers

import (
	"ProjetoUnivesp2020/utils"
)

type sslConfig struct {
	CertPath string
	KeyPath  string
}

type Config struct {
	SSL        *sslConfig
	Bind       string
	BcryptCost int
}

var Configs = LoadConfigs()

func LoadConfigs() *Config {
	return &Config{
		Bind:       utils.GetEnv("BIND", ":1443"),
		BcryptCost: utils.GetIntFromEnv("BCOST", 12),
		SSL: &sslConfig{
			CertPath: utils.GetEnv("CERTPATH", "./certs/server.crt"),
			KeyPath:  utils.GetEnv("KEYPATH", "./certs/server.key"),
		},
	}
}
