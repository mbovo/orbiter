package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"net/http"
)

type AuthConfig struct {
	Enabled bool   `default:"false"`
	User    string `envconfig:"AUTH_USER" default:"orbiter"`
	Pass    string `envconfig:"AUTH_PASS" default:"orbiter"`
	Realm   string `envconfig:"AUTH_REALM" default:"Restricted"`
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
		e := envconfig.Process("ORBITER", &ac)
		if e != nil {
			logrus.Fatal(e.Error())
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="`+ac.Realm+`"`)

		u, p, ok := r.BasicAuth()
		if ok == false {
			w.WriteHeader(401)
			w.Write([]byte("Not Authorized"))
		}

		if ac.User != u || ac.Pass != p {
			w.WriteHeader(401)
			w.Write([]byte("Invalid username or password"))
		}

		h.ServeHTTP(w, r)
	}

}
