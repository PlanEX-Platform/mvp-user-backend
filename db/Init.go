package db

import (
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	//"github.com/spf13/viper"
	"github.com/spf13/viper"
	"time"
)

func Init() *mgo.Session {
	url := viper.GetString("db.url")
	name := viper.GetString("db.name")
	user := viper.GetString("db.user")
	pass := viper.GetString("db.pass")
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{url},
		Timeout:  60 * time.Second,
		Database: name,
		Username: user,
		Password: pass,
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	log.WithError(err).Debug(session)
	err = session.DB(name).Login(user, pass)
	log.WithError(err).Debug(session)
	session.SetMode(mgo.Monotonic, true)
	CreateIndexes(session)
	return session
}
