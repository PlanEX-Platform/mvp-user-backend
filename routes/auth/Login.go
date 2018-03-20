package auth

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	"planex/user-backend/models"
	"planex/user-backend/utils"
)


func Login(basicSession *mgo.Session) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		session := basicSession.Copy()
		defer session.Close()

		// acquiring user input
		r.ParseForm()
		email := r.Form["email"][0]
		password := r.Form["password"][0]

		// TODO: validate inputs

		log.Debug("Trying to login: " + email)

		// checking credentials and sending answer
		if models.NeedConfirmation(email, "register", session) {
			w.Write([]byte("{ status: \"confirmation\" }"))
			log.Debug("Need confirmation: " + email)
			return
		} else {
			user, err := models.UserByEmail(email, session)
			if err == nil {
				if utils.CompareHashAndPass(user.PassHash, password) {
					token, exp := utils.GenJWT(user.ID.String())
					cookie := http.Cookie{Name: "Bearer", Value: token, Expires: exp}
					http.SetCookie(w, &cookie)
					w.Write([]byte("{ status: \"logged\" }"))
					log.Debug("Success login: " + email)
					return
				}
			}
		}

		w.Write([]byte("{ status: \"fail\" }"))
		log.Debug("Wrong password for: " + email)
	}
}
