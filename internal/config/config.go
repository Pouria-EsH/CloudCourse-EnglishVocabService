package config

import "time"

type Config struct {
	Port  string      `koanf:"port"`
	Cache CacheConfig `koanf:"cache"`
	Vocab ExternalAPI `koanf:"english_api"`
}

type CacheConfig struct {
	Address        string        `koanf:"address"`
	ConnectTimeOut time.Duration `koanf:"connection_timeout"`
	CacheTimeOut   time.Duration `koanf:"cache_timeout"`
	Password       string        `koanf:"password"`
}

type ExternalAPI struct {
	ApiKey            string        `koanf:"apikey"`
	ConnectionTimeout time.Duration `koanf:"connection_timeout"`
}
