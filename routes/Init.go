package routes

import (
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/auth0/go-jwt-middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mvp-user-backend/routes/auth"
	"mvp-user-backend/routes/funds"
)

var NotImplemeted = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("{ error: \"Not implemented\" }"))
}

func Init(session *mgo.Session) *httprouter.Router {
	router := httprouter.New()
	router.POST("/api/register", auth.Register(session))
	router.POST("/api/login", auth.Login(session))

	// nonauth -----------------------------------------------------------------------------------------------------------
	serveStatic(router)
	router.POST("/api/confirm", Confirm(session))

	// auth --------------------------------------------------------------------------------------------------------------

	// funds
	router.POST("/api/balances", checkAuth(NotImplemeted))
	router.POST("/api/deposit", checkAuth(funds.Deposit))
	router.POST("/api/withdraw", checkAuth(NotImplemeted))
	router.POST("/api/transfer/history", checkAuth(NotImplemeted))

	// trade
	router.POST("/api/trade/order/make", checkAuth(NotImplemeted))
	router.POST("/api/trade/order/cancel", checkAuth(NotImplemeted))
	router.POST("/api/trade/orders", checkAuth(NotImplemeted))
	router.POST("/api/trade/history", checkAuth(NotImplemeted))

	return router
}

func checkAuth(handleFunc httprouter.Handle) httprouter.Handle {
	var jwtKey = viper.GetString("jwt.secret")
	m := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
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

func serveStatic(router *httprouter.Router) {
	fileServer := http.FileServer(http.Dir("/go/src/mvp-user-backend/static"))
	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		log.Debugf("url: %v", req.URL.Path)
		req.URL.Path = ""
		fileServer.ServeHTTP(w, req)
	})
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Debugf("url: %v", req.URL.Path)
		if !strings.HasPrefix(req.URL.Path, "/static") {
			if !strings.ContainsAny(req.URL.Path, ".") {
				req.URL.Path = ""
			}
		}
		fileServer.ServeHTTP(w, req)
	})
}
