package transport

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/orders", orders).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID:[0-9a-zA-Z]+}", order).Methods(http.MethodGet)
	return logMiddleware(r)
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func orders(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "[{\n    \"id\": \"3fa85f64-5717-4562-b3fc-2c963f66afa6\",\n    \"menuItems\": [{\n        \"id\": \"3fa85f64-5717-4562-b3fc-2c963f66afa6\",\n        \"quantity\": 0\n    }]\n}]")
}

func order(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]
	if id == "3fa85f64-5717-4562-b3fc-2c963f66afa6" {
		fmt.Fprint(w, "{\n  \"id\": \"3fa85f64-5717-4562-b3fc-2c963f66afa6\",\n  \"orderedAtTimestamp\": 0,\n  \"cost\": 0,\n  \"menuItems\": [{\n      \"id\": \"3fa85f64-5717-4562-b3fc-2c963f66afa6\",\n      \"quantity\": 0\n  }]\n}")
	}
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
			"time":       time.Now(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}
