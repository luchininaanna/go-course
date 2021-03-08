package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderserver/pkg/orderservice/transport"
	"os"
	"os/signal"
	"syscall"
)

const serverUrl string = ":8000"
const dbDriver string = "mysql"
const dataSourceName string = "root:1234@/orderservice"

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	killSignalChan := getKillSignalChan()
	srv := startServer(serverUrl)

	createDbConnection(dbDriver, dataSourceName)

	waitForKillSignal(killSignalChan)
	log.Fatal(srv.Shutdown(context.Background()))
}

func startServer(serverUrl string) *http.Server {
	router := transport.Router()
	srv := &http.Server{Addr: serverUrl, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	return srv
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}

func createDbConnection(dbDriver string, dataSourceName string) {
	db, err := sql.Open(dbDriver, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}
