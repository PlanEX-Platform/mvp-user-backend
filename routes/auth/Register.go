package auth

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"planex/user-backend/models"
	"planex/user-backend/mail"
	log "github.com/sirupsen/logrus"
)

func Register(basicSession *mgo.Session) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		session := basicSession.Copy()
		defer session.Close()

		r.ParseForm()
		email := r.Form["email"][0]
		password := r.Form["password"][0]

		// TODO: validate inputs

		if models.CreateUser(email, password, session) {
			if createConfimation(email, session) {
				w.Write([]byte("{ status: \"success\" }"))
				log.Debug("Successful register: " + email)
				return
			} else {
				log.Debug("Failed to create confirmation: " + email)
			}
		} else {
			log.Debug("Failed to create user: " + email)
		}
		w.Write([]byte("{ status: \"fail\" }"))
		log.Debug("Failed register: " + email)
	}
}

func createConfimation(email string, session *mgo.Session) bool {
	created, token := models.CreateConfirmation(email, "register", session)
	if created {
		sended := mail.SendWelcomeLetter(email, token)
		return created && sended
	}
	return false
}