package auth

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	"mvp-user-backend/models"
	"mvp-user-backend/mail"
)

func Register(basicSession *mgo.Session) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		session := basicSession.Copy()
		defer session.Close()

		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")

		// TODO: validate inputs

		if email == "" || password == "" {
			w.Write([]byte(`{ "status": "fail" }`))
			return
		}

		if models.CreateUser(email, password, session) {
			if createConfimation(email, session) {
				w.Write([]byte(`{ "status": "success" }`))
				log.Debug("Successful register: " + email)
				return
			} else {
				log.Debug("Failed to create confirmation: " + email)
			}
		} else {
			log.Debug("Failed to create user: " + email)
		}
		w.Write([]byte(`{ "status": "fail" }`))
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
