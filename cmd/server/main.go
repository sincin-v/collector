package main

import (
	"log"
	"net/http"

	"github.com/sincin-v/collector/internal/server/config"
	"github.com/sincin-v/collector/internal/server/router"
	"github.com/sincin-v/collector/internal/storage"
)

func main() {

	serverConfig, err := config.GetServerConfig()

	if err != nil {
		log.Fatalf("Cannot get server params for start")
	}

	log.Printf("Start server work on %s", serverConfig.Host)

	memStorage := storage.New()

	serverRouter := router.CreateRouter(&memStorage)

	httpErr := http.ListenAndServe(serverConfig.Host, serverRouter)
	if httpErr != nil {
		panic(err)
	}

}
