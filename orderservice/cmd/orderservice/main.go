package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderserver/pkg/orderservice/transport"
	"os"
	"os/signal"
	"syscall"
)

const dataSourcePattern string = "%s:%s@%s/%s?parseTime=true"
const applicationID string = "orderservice"

type config struct {
	ServerPort       string `envconfig:"server_port"`
	DatabaseName     string `envconfig:"database_name"`
	DatabaseAddress  string `envconfig:"database_address"`
	DatabaseUser     string `envconfig:"database_user"`
	DatabasePassword string `envconfig:"database_password"`
	DatabaseDriver   string `envconfig:"database_driver"`
}

func main() {
	conf, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	killSignalChan := getKillSignalChan()
	srv := startServer(*conf)

	waitForKillSignal(killSignalChan)
	log.Fatal(srv.Shutdown(context.Background()))
}

func startServer(conf config) *http.Server {
	dataSourceName := fmt.Sprintf(dataSourcePattern, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseAddress, conf.DatabaseName)
	db := createDbConnection(conf.DatabaseDriver, dataSourceName)
	router := transport.NewRouter(db)
	srv := &http.Server{Addr: ":" + conf.ServerPort, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
		log.Error(db.Close())
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

func createDbConnection(dbDriver string, dataSourceName string) *sql.DB {
	db, err := sql.Open(dbDriver, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func parseConfig() (*config, error) {
	c := config{}
	if err := envconfig.Process(applicationID, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
