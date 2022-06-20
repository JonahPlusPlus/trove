package trove

import (
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/valyala/fasthttp"
)

type Producers struct {
	requestEventProducer sarama.AsyncProducer
}

func newProducers(brokers []string) Producers {
	config := sarama.NewConfig()
	config.ClientID = "trove"

	requestEventProducer, err := sarama.NewAsyncProducer(brokers, config)

	if err != nil {
		log.Fatal(err)
	}

	return Producers{
		requestEventProducer,
	}
}

func (p *Producers) logRequest(ctx *fasthttp.RequestCtx, duration float64) {
	event := &RequestEvent{
		Method: string(ctx.Method()),
		Host:   ctx.LocalAddr().String(),
		Path:   string(ctx.Path()),
		Remote: ctx.RemoteAddr().String(),
		Time:   duration,
	}

	p.requestEventProducer.Input() <- &sarama.ProducerMessage{
		Topic: "requests",
		Key:   sarama.ByteEncoder(ctx.RemoteIP()),
		Value: event,
	}
}

func init() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
}
