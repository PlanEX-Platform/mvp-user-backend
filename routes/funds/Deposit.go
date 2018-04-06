package funds

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"io"
)

func Deposit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if user := (r.Context().Value("user")); user != nil {
		id := (user.(*jwt.Token)).Claims.(jwt.MapClaims)["id"].(string)
		log.Infof("Deposit request user_id: %v", id)
		var url = viper.GetString("deposit_service.url")
		resp, err := http.Get(url + "/account/" + id)
		log.WithError(err).Debugf("Call to deposit service result: %v", resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()
		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}