package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

const (
	prefix    = "ENGLISHVOCAB_"
	delimeter = "."
	seprator  = "__"
)

func callbackEnv(source string) string {
	base := strings.ToLower(strings.TrimPrefix(source, prefix))
	return strings.ReplaceAll(base, seprator, delimeter)
}

func New() Config {
	k := koanf.New(".")

	if err := k.Load(structs.Provider(defaultConfig(), "koanf"), nil); err != nil {
		log.Printf("error loading config.toml: %s", err)
	}

	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Printf("error loading config.toml: %s", err)
	}

	if err := k.Load(env.Provider(prefix, delimeter, callbackEnv), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	var conf Config
	if err := k.Unmarshal("", &conf); err != nil {
		log.Fatalf("error unmarshaling config: %s", err)
	}

	log.Printf("%+v", conf)

	return conf
}
