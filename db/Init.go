package db

import (
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	//"github.com/spf13/viper"
	"github.com/spf13/viper"
)

func Init() *mgo.Session {
	url := viper.GetString("db.url")
	name := viper.GetString("db.name")
	user := viper.GetString("db.user")
	pass := viper.GetString("db.pass")
	session, err := mgo.Dial(url)
	log.WithError(err).Debug(session)
	err = session.DB(name).Login(user, pass)
	log.WithError(err).Debug(session)
	CreateIndexes(session)
	return session
}
