package main

import (
	"log"
	"net/http"

	"github.com/sincin-v/collector/internal/server/config"
	"github.com/sincin-v/collector/internal/server/router"
	"github.com/sincin-v/collector/internal/storage"
)

func main() {

	serverConfig := config.GetServerConfig()

	log.Printf("Start server work on %s", serverConfig.Host)

	metricStorage := storage.New()

	serverRouter := router.CreateRouter(&metricStorage)

	err := http.ListenAndServe(serverConfig.Host, serverRouter)
	if err != nil {
		panic(err)
	}

}
