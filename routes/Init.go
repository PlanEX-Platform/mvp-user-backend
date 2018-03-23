package routes

import (
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/auth0/go-jwt-middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mvp-user-backend/routes/auth"
)

var jwtKey = viper.GetString("jwt.secret")
var NotImplemeted = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("{ error: \"Not implemented\" }"))
}

func Init(session *mgo.Session) *httprouter.Router {
	router := httprouter.New()
	router.POST("/register", auth.Register(session))
	router.POST("/login", auth.Login(session))

	// nonauth -----------------------------------------------------------------------------------------------------------
	router.GET("/", NotImplemeted)
	router.POST("/confirm", Confirm(session))

	// auth --------------------------------------------------------------------------------------------------------------

	// funds
	router.POST("/balances", checkAuth(NotImplemeted))
	router.POST("/deposit", checkAuth(NotImplemeted))
	router.POST("/withdraw", checkAuth(NotImplemeted))
	router.POST("/transfer/history", checkAuth(NotImplemeted))

	// trade
	router.POST("/trade/order/make", checkAuth(NotImplemeted))
	router.POST("/trade/order/cancel", checkAuth(NotImplemeted))
	router.POST("/trade/orders", checkAuth(NotImplemeted))
	router.POST("/trade/history", checkAuth(NotImplemeted))

	return router
}

func checkAuth(handleFunc httprouter.Handle) httprouter.Handle {
	m := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := m.CheckJWT(w, r)
		if err == nil {
			handleFunc(w, r, ps)
		} else {
			log.WithError(err).Debug("JWT authorization failed")
		}
	}
}
