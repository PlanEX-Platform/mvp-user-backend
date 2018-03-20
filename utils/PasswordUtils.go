package utils

import (
	"golang.org/x/crypto/bcrypt"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func CompareHashAndPass(hash string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	log.WithError(err).Debug("Valid password: " + strconv.FormatBool(err == nil))
	return err == nil
}