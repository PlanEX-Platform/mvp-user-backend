package models

import (
	"gopkg.in/mgo.v2/bson"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"golang.org/x/crypto/bcrypt"
	"github.com/spf13/viper"
)

type User struct {
	Id bson.ObjectId `bson:"_id" json:"id"`
	Email	string 			`json:"email" bson:"email"`
	PassHash string 	`json:"pass_hash" bson:"pass_hash"`
}

func CreateUser(email string, pass string, session *mgo.Session) bool {
	dbName := viper.GetString("db.name")
	userCollectionName := viper.GetString("db.user_collection")
	users := session.DB(dbName).C(userCollectionName)
	user := &User{
		Id: bson.NewObjectId(),
		Email: email,
		PassHash: hashAndSalt(pass) }
	err := users.Insert(user)
	if mgo.IsDup(err) {
		log.Debugf("User already exist: %v", email)
	}
	log.WithError(err).Debugf("Trying to create new user: %v", email)
	return err == nil
}

func UserByEmail(email string, session *mgo.Session) (User, error) {
	dbName := viper.GetString("db.name")
	userCollectionName := viper.GetString("db.user_collection")
	users := session.DB(dbName).C(userCollectionName)
	var user User
	err := users.Find(bson.M{"email": email}).One(&user)
	log.WithError(err).Debug(user)
	return user, err
}

func hashAndSalt(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	log.WithError(err).Debug("Password hash generated: " + string(hash))
	return string(hash)
}
