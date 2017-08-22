package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"net/http"
)

type AuthConfig struct {
	AuthEnabled bool   `split_words:"true" default:"false"`
	AuthUser    string `split_words:"true" default:"orbiter"`
	AuthPass    string `split_words:"true" default:"orbiter"`
	AuthRealm   string `split_words:"true" default:"Restricted"`
}

func wrap(h http.HandlerFunc, funx ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, f := range funx {
		h = f(h)
	}
	return h
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var ac AuthConfig
		e := envconfig.Process("orbiter", &ac)
		if e != nil {
			logrus.Fatal(e.Error())
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="`+ac.AuthRealm+`"`)

		u, p, ok := r.BasicAuth()
		if ok == false {
			logrus.Error("No Authentication supplied")
			w.WriteHeader(401)
			w.Write([]byte("Not Authorized"))
			return
		}

		if ac.AuthUser != u || ac.AuthPass != p {
			logrus.Warnf("Invalid username or password for user %s", u)
			w.WriteHeader(401)
			w.Write([]byte("Invalid username or password"))
			return
		}

		logrus.Infof("Succesfully authenticated with user: %s", u)

		h.ServeHTTP(w, r)
	}

}
