package main

import (
	"log"
	"net/http"
)

const webPort = "8081"

type Config struct{}

func main() {
	app := Config{}
	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server

	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	// start http server
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
