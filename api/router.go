package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/gianarb/orbiter/core"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

func GetRouter(core *core.Core, eventChannel chan *logrus.Entry) *mux.Router {
	r := mux.NewRouter()

	var ac AuthConfig
	envconfig.Usage("orbiter", &ac)

	e := envconfig.Process("orbiter", &ac)
	if e != nil {
		logrus.Fatal(e.Error())
	}

	logrus.Infof("Env: %v", ac)

	if ac.AuthEnabled {
		logrus.Info("Enabling Authentication")
		r.HandleFunc("/v1/orbiter/handle/{autoscaler_name}/{service_name}", wrap(Handle(&core.Autoscalers), basicAuth)).Methods("POST")
		r.HandleFunc("/v1/orbiter/handle/{autoscaler_name}/{service_name}/{direction}", wrap(Handle(&core.Autoscalers), basicAuth)).Methods("POST")
		r.HandleFunc("/v1/orbiter/autoscaler", AutoscalerList(core.Autoscalers)).Methods("GET")
		r.HandleFunc("/v1/orbiter/health", Health()).Methods("GET")
		r.HandleFunc("/v1/orbiter/events", Events(eventChannel)).Methods("GET")

	} else {

		r.HandleFunc("/v1/orbiter/handle/{autoscaler_name}/{service_name}", Handle(&core.Autoscalers)).Methods("POST")
		r.HandleFunc("/v1/orbiter/handle/{autoscaler_name}/{service_name}/{direction}", Handle(&core.Autoscalers)).Methods("POST")
		r.HandleFunc("/v1/orbiter/autoscaler", AutoscalerList(core.Autoscalers)).Methods("GET")
		r.HandleFunc("/v1/orbiter/health", Health()).Methods("GET")
		r.HandleFunc("/v1/orbiter/events", Events(eventChannel)).Methods("GET")
	}

	// This lines will be removed October 2017. They are here to offer a soft migation path.
	r.HandleFunc("/handle/{autoscaler_name}/{service_name}", Handle(&core.Autoscalers)).Methods("POST")
	r.HandleFunc("/handle/{autoscaler_name}/{service_name}/{direction}", Handle(&core.Autoscalers)).Methods("POST")
	r.HandleFunc("/autoscaler", AutoscalerList(core.Autoscalers)).Methods("GET")
	r.HandleFunc("/health", Health()).Methods("GET")
	r.HandleFunc("/events", Events(eventChannel)).Methods("GET")

	r.NotFoundHandler = NotFound{}
	return r
}
