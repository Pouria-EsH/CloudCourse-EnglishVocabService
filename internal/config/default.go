package config

import "time"

func defaultConfig() Config {
	return Config{
		Port: "8080",
		Cache: CacheConfig{
			ConnectTimeOut: time.Second,
			CacheTimeOut:   time.Minute,
			Password:       "admin",
		},
		Vocab: ExternalAPI{
			ApiKey:            "",
			ConnectionTimeout: time.Second,
		},
	}
}
