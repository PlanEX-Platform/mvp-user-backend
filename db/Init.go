package db

import (
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	//"github.com/spf13/viper"
	"github.com/spf13/viper"
)

func Init() *mgo.Session {
	url := viper.GetString("db.url")
	session, err := mgo.Dial(url)
	log.WithError(err).Debug(session)
	CreateIndexes(session)
	return session
}
