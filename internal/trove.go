package trove

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/JonahPlusPlus/trove/templates"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Trove struct {
	config       Config
	server       fasthttp.Server
	producers    Producers
	consumers    Consumers
	mongo_client *mongo.Client
	ext_match    *regexp.Regexp
}

func New(config ...Config) Trove {
	c := getConfig(config...)

	matcher, err := regexp.Compile(`\..*`)

	if err != nil {
		log.Fatal(err)
	}

	t := Trove{
		config:       c,
		producers:    newProducers(c.Brokers),
		consumers:    newConsumers(c.Brokers, c.GroupID),
		mongo_client: nil,
		ext_match:    matcher,
	}

	t.server = fasthttp.Server{
		Name:                 "trove",
		Handler:              t.handler,
		ReadTimeout:          5 * time.Second,
		WriteTimeout:         10 * time.Second,
		MaxConnsPerIP:        500,
		MaxRequestsPerConn:   500,
		MaxKeepaliveDuration: 5 * time.Second,
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connecting to MongoDB at", t.config.MongoURI)
	opts := options.Client()
	opts = opts.ApplyURI(t.config.MongoURI)

	t.mongo_client, err = mongo.Connect(context.Background(), opts)

	if err != nil {
		log.Fatal(err)
	}

	if err = t.mongo_client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	return t
}

func (t *Trove) handler(ctx *fasthttp.RequestCtx) {
	log.Println(ctx)

	started := time.Now()
	path := string(ctx.Path())
	var err error

	if strings.Contains(path, ".") {
		// Handle static files
		var f []byte
		f, err = ioutil.ReadFile("./static" + path)
		if err != nil {
			log.Println(err)
		} else {
			ctx.SetBody(f)
			m := mime.TypeByExtension(t.ext_match.FindString(path))
			ctx.SetContentType(m)
		}
	} else {
		// Handle web pages
		ctx.SetContentType("text/html")
		switch path {
		case "/":
			ctx.SetBody([]byte(templates.PrintPage(&templates.RootPage{})))
		case "/inventory":
			col := t.mongo_client.Database("trove").Collection("items")
			log.Println(col.Name())

			projection := bson.D{{"name", 1}, {"quantity", 1}, {"_id", 0}}
			opts := options.Find().SetProjection(projection)

			cur, err := col.Find(context.Background(), bson.D{}, opts)

			if err != nil {
				log.Println(err)
				break
			}

			var results []bson.D
			if err = cur.All(context.Background(), &results); err != nil {
				log.Println(err)
				break
			}

			ctx.SetBody([]byte(templates.PrintPage(&templates.InventoryPage{Items: results})))
		case "/dashboard":
			ctx.SetBody([]byte(templates.PrintPage(&templates.DashboardPage{})))
		case "/ws":
			t.consumers.analytics.ws_handler(ctx)
		default:
			ctx.SetBody([]byte(templates.PrintPage(&templates.UnfoundPage{})))
			err = errors.New("404")
		}
	}

	t.producers.logRequest(ctx, float64(time.Since(started))/float64(time.Second), err)
}

func (t *Trove) Run() {
	t.consumers.exitCallback = t.consumers.consume()

	ln, err := net.Listen("tcp", t.config.Address)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Serving TLS connections on " + t.config.Address)
	if err = t.server.ServeTLS(ln, t.config.Certificate, t.config.Key); err != nil {
		log.Fatal(err)
	}
}

func (t *Trove) Exit() {
	if t.consumers.exitCallback != nil {
		t.consumers.exitCallback()
	}
	t.server.Shutdown()
	if t.mongo_client != nil {
		t.mongo_client.Disconnect(context.Background())
	}
}
