package main

import (
	log "github.com/sirupsen/logrus"
	"planex/user-backend/db"
	//"planex/user-backend/routes"
	"planex/user-backend/logenv"
	"planex/user-backend/config"
)

func init() {
	logenv.InitLog()
	config.Load()
}

func main() {
	session := db.Init()
	log.Info(session)
	//routes.Init(session)
	// broadcasting.Start()
	// market.Start()
}
