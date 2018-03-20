package db

import (
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func CreateIndexes(s *mgo.Session) {
	createUserIndex(s)
}

func createUserIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	dbName := viper.GetString("db.name")
	userCollectionName := viper.GetString("db.user_collection")
	users := session.DB(dbName).C(userCollectionName)

	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	}
	err := users.EnsureIndex(index)
	if err != nil {
		log.Panic(err)
	}
}
