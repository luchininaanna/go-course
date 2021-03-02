package transport

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"orderserver/model"
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
	orders := model.OrderList{
		Orders: []model.Order{
			{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", MenuItems: []model.MenuItem{{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Quantity: 0}}},
		},
	}

	jsonOrders, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, string(jsonOrders)); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}

func order(w http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["ID"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprint(w, "Order not found")
		if err != nil {
			log.Error(err)
		}
		return
	}

	detailedOrder := model.DetailedOrder{
		Order: model.Order{ID: id, MenuItems: []model.MenuItem{{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Quantity: 0}}},
		Cost:  1,
		Time:  1,
	}

	jsonDetailedOrder, err := json.Marshal(detailedOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, string(jsonDetailedOrder)); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
			"time":       start,
			"reqDur":     time.Now().Sub(start),
		}).Info("got a new request")

	})
}
