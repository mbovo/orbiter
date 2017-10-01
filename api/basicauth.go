package api

import (
	"github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type AuthConfig struct {
	AuthEnabled     bool              `split_words:"true" default:"false"`
	AuthRealm       string            `split_words:"true" default:"Restricted"`
	AuthCredentials map[string]string `split_words:"true" default:"orbiter:orbiter"`
}

func wrap(h http.HandlerFunc, funx ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, f := range funx {
		h = f(h)
	}
	return h
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ac := new(AuthConfig)

		s, e := os.LookupEnv("ORBITER_AUTH_ENABLED")
		ac.AuthEnabled, _ = strconv.ParseBool(s)
		if !e {
			ac.AuthEnabled = e
		}
		s, e = os.LookupEnv("ORBITER_AUTH_REALM")
		ac.AuthRealm = s
		if !e {
			ac.AuthRealm = "Restricted"
		}
		s, e = os.LookupEnv("ORBITER_AUTH_CREDENTIALS")
		if e {
			for _, pair := range strings.Split(s, ",") {
				auth := strings.Split(pair, ":")
				ac.AuthCredentials[auth[0]] = auth[1]
			}
		} else {
			ac.AuthCredentials["orbiter"] = "orbiter"
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="`+ac.AuthRealm+`"`)

		u, p, ok := r.BasicAuth()
		if ok == false {
			logrus.Error("No Authentication supplied")
			w.WriteHeader(401)
			w.Write([]byte("Not Authorized"))
			return
		}

		if ac.AuthCredentials[u] != p {
			logrus.Warnf("Invalid username or password for user %s", u)
			w.WriteHeader(401)
			w.Write([]byte("Invalid username or password"))
			return
		}

		logrus.Infof("Succesfully authenticated with user: %s", u)

		h.ServeHTTP(w, r)
	}

}
