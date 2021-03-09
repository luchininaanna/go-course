package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"orderserver/pkg/orderservice/repository"
	"orderserver/pkg/orderservice/service"
	"time"
)

type server struct {
	orderService service.OrderService
}

func newServer(db *sql.DB) *server {
	r := repository.NewOrderRepository(db)
	s := service.NewOrderService(r)
	return &server{s}
}

func NewRouter(db *sql.DB) http.Handler {
	srv := newServer(db)
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/orders", getOrders).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID:[0-9a-zA-Z]+}", getOrderInfo).Methods(http.MethodGet)
	s.HandleFunc("/order", srv.createOrder).Methods(http.MethodPost)
	return logMiddleware(r)
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, "Hello World!")
	if err != nil {
		log.Error(err)
	}
}

func getOrders(w http.ResponseWriter, _ *http.Request) {
	orders := service.OrderList{
		Orders: []service.Order{
			{MenuItems: []service.MenuItem{{ID: uuid.New().String(), Quantity: 0}}},
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getOrderInfo(w http.ResponseWriter, r *http.Request) {
	_, found := mux.Vars(r)["ID"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprint(w, "Order not found")
		if err != nil {
			log.Error(err)
		}
		return
	}

	detailedOrder := service.DetailedOrder{
		Order: service.Order{MenuItems: []service.MenuItem{{ID: uuid.New().String(), Quantity: 0}}},
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (srv *server) createOrder(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Can't read request body with error %v", err)
	}

	defer r.Body.Close()

	var orderData service.Order
	err = json.Unmarshal(b, &orderData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Can't parse json response with error %v", err)
	}

	err = srv.orderService.AddOrder(orderData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
			"reqDur":     time.Since(start).String(),
		}).Info("got a new request")
	})
}
