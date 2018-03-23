package models

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mvp-user-backend/utils"
)

type Confirm struct {
	ID bson.ObjectId 	`json:"_id" bson:"_id"`
	Email string 			`json:"email" bson:"email"`
	Token string			`json:"token" bson:"token"`
	Category string				`json:"category" bson:"category"`
}

var confirmCollectionName = viper.GetString("db.confirm_collection")

func CreateConfirmation(email string, category string, session *mgo.Session) (bool, string) {
	confirmations := session.DB(dbName).C(confirmCollectionName)
	token := utils.GenConfirmationToken()
	confirm := Confirm{
		Email: email,
		Category: category }
	err := confirmations.Insert(confirm)
	log.WithError(err).Debugf("Trying to create confirmation instance: %v category: %v", email, category)
	return err == nil, token
}

func NeedConfirmation(email string, category string, session *mgo.Session) bool {
	confirmations := session.DB(dbName).C(confirmCollectionName)
	var confirm Confirm
	err := confirmations.Find(bson.M{"email": email, "category": category}).One(&confirm)
	log.WithError(err).Debug(confirm)
	return err == nil
}

func RemoveConfirmation(token string, category string, session *mgo.Session) (bool, string) {
	confirmations := session.DB(dbName).C(confirmCollectionName)
	var confirm Confirm
	confirmations.Find(bson.M{"token": token, "category": category}).One(&confirm)
	err := confirmations.Remove(bson.M{"token": token})
	log.WithError(err).Debugf("Removing confirmation by token: %v", token)
	return err == nil, confirm.Email
}