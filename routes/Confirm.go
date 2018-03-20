package routes

import (
	"gopkg.in/mgo.v2"
	"planex/user-backend/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"planex/user-backend/mail"
)

func Confirm(basicSession *mgo.Session) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		session := basicSession.Copy()
		defer session.Close()

		// acquiring confirmation token
		r.ParseForm()
		token := r.Form["token"][0]
		category := r.Form["cat"][0]

		// TODO: validate inputs

		// trying to confirm and sending answer
		if confirmByToken(token, category, session) {
			w.Write([]byte("{ status: \"confirmed\" }"))
			log.Debugf("Successfully confirmed %v by %v", category, token)
		} else {
			w.Write([]byte("{ status: \"fail\" }"))
			log.Debugf("Failed confirmation %v by %v", category, token)
		}
	}
}

func confirmByToken(token string, category string, session *mgo.Session) bool {
	removed, email := models.RemoveConfirmation(token, category, session)
	if email != "" {
		sended := mail.SendSuccessConfirmation(email, category)
		return removed && sended
	}
	return false
}
