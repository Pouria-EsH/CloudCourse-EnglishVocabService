package service

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"vocabsrv/internal/cache"
	"vocabsrv/internal/monitor"
	"vocabsrv/internal/vocab"

	"github.com/labstack/echo/v4"
)

type VacabService struct {
	port        string
	vocabClient vocab.ApiNinjas
	cache       cache.RedisCache
	metr        monitor.PromMetrics
}

type DefinitionResponse struct {
	FetchStatus string `json:"Fetch-Status"`
	Definition  string `json:"Definition"`
}

func NewVacabService(port string, voc vocab.ApiNinjas,
	cachedb cache.RedisCache, met monitor.PromMetrics) *VacabService {
	return &VacabService{
		port:        port,
		vocabClient: voc,
		cache:       cachedb,
		metr:        met,
	}
}

func (vs VacabService) DefinitionRequestHandler(e echo.Context) error {
	wrd := e.QueryParam("word")
	if wrd == "" {
		log.Println("no word given")
		return e.String(http.StatusBadRequest, "no word given")
	}

	fetchstatus := "from-cache"

	definition, err := vs.cache.GetWord(wrd)
	if err == cache.ErrNotCached {
		fetchstatus = "from-api"
		definition, err = vs.vocabClient.GetDefinition(wrd)
		if err != nil {
			log.Println("error at vocab client GetDefinition: ", err)
			return err
		}

		err = vs.cache.AddWord(wrd, definition)
		if err != nil {
			log.Println("error at word caching: ", err)
			return err
		}
	} else if err != nil {
		log.Println("error at reading cache: ", err)
		return err
	} else {
		vs.metr.RedisCounter.Inc()
	}

	return e.JSON(http.StatusOK, DefinitionResponse{
		FetchStatus: fetchstatus,
		Definition:  definition,
	})
}

func (vs VacabService) RandomWrdRequestHandler(e echo.Context) error {
	word, err := vs.vocabClient.GetRandom()
	if err != nil {
		fmt.Println("error at vocab client GetRandom: ", err)
		return err
	}

	return e.String(http.StatusOK, word)
}

func (vs VacabService) Execute() error {
	app := echo.New()
	app.Use(
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if c.Path() == "/definition" {
					vs.metr.DefinitionRequestCount.Inc()
				} else if c.Path() == "/randword" {
					vs.metr.RandwordRequestCount.Inc()
				}
				start := time.Now()
				ret := next(c)
				vs.metr.ApiLatency.Observe(time.Since(start).Seconds())
				return ret
			}
		},
	)
	app.GET("/definition", vs.DefinitionRequestHandler)
	app.GET("/randword", vs.RandomWrdRequestHandler)
	return app.Start(":" + vs.port)
}
