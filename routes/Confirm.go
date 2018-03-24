package routes

import (
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"mvp-user-backend/models"
	"mvp-user-backend/mail"
)

func Confirm(basicSession *mgo.Session) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		session := basicSession.Copy()
		defer session.Close()

		// acquiring confirmation token
		r.ParseForm()
		token := r.FormValue("token")
		category := r.FormValue("cat")

		// TODO: validate inputs

		if token == "" || category == "" {
			w.Write([]byte(`{ "status": "fail" }`))
			return
		}

		// trying to confirm and sending answer
		if confirmByToken(token, category, session) {
			w.Write([]byte(`{ "status": "confirmed" }`))
			log.Debugf("Successfully confirmed %v by %v", category, token)
		} else {
			w.Write([]byte(`{ "status": "fail" }`))
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
