package main

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	trove "github.com/JonahPlusPlus/trove/internal"
)

var config trove.Config

func init() {
	content, err := ioutil.ReadFile("trove.toml")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = toml.Decode(string(content), &config)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	log.Printf("Kafka Broker: %s", config.Broker)

	trove := trove.New(config)

	trove.Run()
}
