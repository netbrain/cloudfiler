package main

import (
	"github.com/netbrain/cloudfiler/app"
	"github.com/netbrain/cloudfiler/app/conf"
	"log"
	"net/http"
)

func main() {
	serverAddr := conf.Config.ServerAddr
	log.Printf("Starting server, listening on: %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, app.WebHandler); err != nil {
		panic(err)
	}
}
