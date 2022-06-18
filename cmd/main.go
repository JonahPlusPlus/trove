package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	trove "github.com/JonahPlusPlus/trove/pkg"
	"github.com/Shopify/sarama"
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
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	log.Printf("Kafka Broker: %s", config.Broker)

}
