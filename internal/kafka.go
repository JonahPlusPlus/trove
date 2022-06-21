package trove

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

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

func (p *Producers) logRequest(ctx *fasthttp.RequestCtx, duration float64, err error) {

	var e *string = nil

	if err != nil {
		msg := err.Error()
		e = &msg
	}

	event := &RequestEvent{
		Method: string(ctx.Method()),
		Host:   ctx.LocalAddr().String(),
		Path:   string(ctx.Path()),
		Remote: ctx.RemoteAddr().String(),
		Time:   duration,
		Error:  e,
	}

	p.requestEventProducer.Input() <- &sarama.ProducerMessage{
		Topic: "requests",
		Key:   sarama.ByteEncoder(ctx.RemoteIP()),
		Value: event,
	}
}

type Consumers struct {
	requestEventConsumer sarama.ConsumerGroup
	analytics            Analytics
	exitCallback         func()
}

func newConsumers(brokers []string, groupID string) Consumers {
	config := sarama.NewConfig()
	config.ClientID = "trove-analytics"

	requestEventConsumer, err := sarama.NewConsumerGroup(brokers, groupID, config)

	if err != nil {
		log.Fatal(err)
	}

	return Consumers{
		requestEventConsumer: requestEventConsumer,
		analytics:            newAnalytics(),
		exitCallback:         nil,
	}
}

func (c *Consumers) consume() func() {
	consumer := RequestEventConsumer{
		ready:     make(chan bool),
		analytics: &c.analytics,
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := c.requestEventConsumer.Consume(ctx, []string{"requests"}, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	log.Println("Sarama consumer up and running!...")

	return func() {
		log.Println("Closing consumer")
		cancel()
		wg.Wait()
		if err := c.requestEventConsumer.Close(); err != nil {
			log.Panicf("Error closing consumer: %v", err)
		}
	}
}

type RequestEventConsumer struct {
	ready     chan bool
	analytics *Analytics
}

func (consumer *RequestEventConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

func (consumer *RequestEventConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *RequestEventConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")

			var request RequestEvent

			err := json.Unmarshal(message.Value, &request)

			if err != nil {
				log.Println("Can't unmarshal message")
			}

			consumer.analytics.addRequestEvent(request)

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func init() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
}
