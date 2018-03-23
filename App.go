package main

import (
	log "github.com/sirupsen/logrus"
	"mvp-user-backend/config"
	"mvp-user-backend/routes"
	"mvp-user-backend/logenv"
	"mvp-user-backend/db"
	"net/http"
)

func init() {
	logenv.InitLog()
	config.Load()
}

func main() {
	session := db.Init()
	log.Info(session)
	router := routes.Init(session)
	log.Debug("Starting at 7200...")
	http.ListenAndServe(":7200", router)
	// broadcasting.Start()
	// market.Start()
}
