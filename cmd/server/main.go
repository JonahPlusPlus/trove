package main

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	trove "github.com/JonahPlusPlus/trove/internal"
	"github.com/valyala/fasthttp"
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

	log.Fatal(fasthttp.ListenAndServeTLS(trove.Address(), trove.CertificatePath(), trove.KeyPath(), trove.Run))
}
