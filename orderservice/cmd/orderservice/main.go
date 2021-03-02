package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"orderserver/pkg/orderservice"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	serverUrl := ":8000"

	killSignalChan := getKillSignalChan()
	srv := startServer(serverUrl)

	waitForKillSignal(killSignalChan)
	log.Fatal(srv.Shutdown(context.Background()))
}

func startServer(serverUrl string) *http.Server {
	router := orderservice.Router()
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
