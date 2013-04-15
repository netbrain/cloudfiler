package main

import (
	"github.com/netbrain/cloudfiler/app"
	"github.com/netbrain/cloudfiler/app/conf"
	"log"
	"net/http"
)

func main() {
	log.Print("Creating static files handler")
	staticPfx := "/static/"
	staticFilesPath := conf.Config.ApplicationHome + staticPfx
	staticHandler := http.StripPrefix(staticPfx, http.FileServer(http.Dir(staticFilesPath)))
	http.Handle(staticPfx, staticHandler)

	log.Print("Creating application handler")
	http.Handle("/", app.WebHandler)

	serverAddr := conf.Config.ServerAddr
	log.Printf("Starting server, listening on: %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		panic(err)
	}
}
