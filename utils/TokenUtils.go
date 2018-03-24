package utils

import (
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"time"
	"crypto/rand"
	"encoding/base64"
	"github.com/spf13/viper"
)

func GenJWT(id string) (string, time.Time) {
	jwtKey := viper.GetString("jwt.secret")
	expDelta := viper.GetInt64("jwt.expiration_delta")
	token := jwt.New(jwt.SigningMethodHS256)

	expiration := time.Now().Add(time.Hour * time.Duration(expDelta))
	token.Claims = jwt.MapClaims{
		"exp": expiration.Unix(),
		"iat": time.Now().Unix(),
		"id":  id }

	tokenString, err := token.SignedString([]byte(jwtKey))
	log.WithError(err).Debug("JWT token generated: " + tokenString)

	return tokenString, expiration
}

func GenConfirmationToken() string {
	size := 32
	rb := make([]byte,size)
	_, err := rand.Read(rb)
	if err != nil {
		log.Error(err)
	}
	rs := base64.URLEncoding.EncodeToString(rb)
	return rs
}
