package main

import (
	log "github.com/sirupsen/logrus"
	"mvp-user-backend/config"
	"mvp-user-backend/routes"
	"mvp-user-backend/logenv"
	"mvp-user-backend/db"
)

func init() {
	logenv.InitLog()
	config.Load()
}

func main() {
	session := db.Init()
	log.Info(session)
	routes.Init(session)
	// broadcasting.Start()
	// market.Start()
}
