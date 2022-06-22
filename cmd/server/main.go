package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	trove "github.com/JonahPlusPlus/trove/internal"
)

var config trove.Config

func init() {
	content, err := ioutil.ReadFile(".config")
	path := string(content)
	if err != nil {
		log.Fatal(err.Error())
	}
	content, err = ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = toml.Decode(string(content), &config)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	log.Printf("Kafka Broker: %s", config.Brokers)

	trove := trove.New(config)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Received an interrupt, exiting...")
		trove.Exit()
		os.Exit(0)
	}()

	trove.Run()
}
