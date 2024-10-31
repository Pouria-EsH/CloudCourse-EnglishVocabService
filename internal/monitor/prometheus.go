package monitor

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetrics struct {
	RandwordRequestCount   prometheus.Counter
	DefinitionRequestCount prometheus.Counter
	RedisCounter           prometheus.Counter
	// apiLatency			prometheus.
}

func NewPromMetrics() *PromMetrics {
	m := PromMetrics{
		RandwordRequestCount: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "randomWordTotalRequest",
				Help: "Total number of requests to /ranword api",
			},
		),
		DefinitionRequestCount: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "definitionWordTotalRequest",
				Help: "Total number of requests to /deinition api",
			},
		),
		RedisCounter: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "redisTotalRequest",
				Help: "Total number of requests answered by redis",
			},
		),
	}

	prometheus.MustRegister(
		m.RandwordRequestCount,
		m.DefinitionRequestCount,
		m.RedisCounter,
	)

	return &m
}

func (pm PromMetrics) Expose(endpoint string, port uint16) {
	p := strconv.FormatUint(uint64(port), 10)
	http.Handle(endpoint, promhttp.Handler())
	err := http.ListenAndServe(":"+p, nil)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Printf("prometheus Failed To Start : %v", err)
	}
}
