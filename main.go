package main

import (
	"log"
	"vocabsrv/internal/cache"
	"vocabsrv/internal/config"
	"vocabsrv/internal/monitor"
	"vocabsrv/internal/service"
	"vocabsrv/internal/vocab"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("credentials.env"); err != nil {
		log.Println("no 'credentials.env' file found")
	}

	cfg := config.New()

	redisdb := cache.NewRedisClient(
		cfg.Cache.Address,
		cfg.Cache.Password,
		cfg.Cache.CacheTimeOut,
		cfg.Cache.ConnectTimeOut,
	)

	ninjas_client := vocab.NewApiNinjas(
		cfg.Vocab.ApiKey,
		cfg.Vocab.ConnectionTimeout,
	)

	metrs := monitor.NewPromMetrics()
	go metrs.Expose("/metrics", 9091)

	al := service.NewVacabService(cfg.Port, *ninjas_client, *redisdb, *metrs)

	err := al.Execute()
	if err != nil {
		log.Fatalf("could not start service1: %v\n", err)
	}
}
