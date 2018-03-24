package auth

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	"mvp-user-backend/models"
	"mvp-user-backend/utils"
)


func Login(basicSession *mgo.Session) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		session := basicSession.Copy()
		defer session.Close()

		// acquiring user input
		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")

		// TODO: validate inputs

		log.Debug("Trying to login: " + email)

		if email == "" || password == "" {
			w.Write([]byte(`{ "status": "fail" }`))
			return
		}

		// checking credentials and sending answer
		if models.NeedConfirmation(email, "register", session) {
			w.Write([]byte(`{ "status": "confirmation" }`))
			log.Debug("Need confirmation: " + email)
			return
		} else {
			user, err := models.UserByEmail(email, session)
			if err == nil {
				if utils.CompareHashAndPass(user.PassHash, password) {
					token, exp := utils.GenJWT(user.Id.String())
					log.Debugf("token: %v exp: %v", token, exp)
					// cookie := http.Cookie{
					// 	Name: "Bearer",
					// 	Value: token,
					// 	Expires: exp,
					// 	HttpOnly: true}
					// log.Debugf("Login cookie: %v", cookie)
					// http.SetCookie(w, &cookie)
					w.Write([]byte(`{ "status": "logged", "token": "` + token + `" }`))
					log.Debug("Success login: " + email)
					return
				}
			}
		}

		w.Write([]byte(`{ "status": "fail" }`))
		log.Debug("Wrong password for: " + email)
	}
}
