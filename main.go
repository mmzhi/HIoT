package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"

	_ "github.com/fhmq/hmq/adapter/impl"
	"github.com/fhmq/hmq/broker"
	_ "github.com/fhmq/hmq/database/mysql"
	_ "github.com/fhmq/hmq/database/sqlite"
	_ "github.com/fhmq/hmq/plugins/manage/impl"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config, err := broker.ConfigureConfig()
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}

	b, err := broker.NewBroker(config)
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
