package main

import (
	"github.com/fhmq/hmq/config"
	"github.com/fhmq/hmq/mqtt/broker"
	"log"
	"os"
	"os/signal"
	"runtime"

	_ "github.com/fhmq/hmq/database/mysql"
	_ "github.com/fhmq/hmq/database/sqlite"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg, err := config.Configure()
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}

	b, err := broker.NewBroker(cfg)
	if err != nil {
		log.Fatal("New Broker error: ", err)
	}
	b.Start()

	s := waitForSignal()
	log.Println("signal received, broker closed.", s)
}

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}
