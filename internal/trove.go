package trove

import (
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
)

type Trove struct {
	config    Config
	server    fasthttp.Server
	producers Producers
	consumers Consumers
	ext_match *regexp.Regexp
}

func New(config ...Config) Trove {
	c := getConfig(config...)

	matcher, err := regexp.Compile(`\..*`)

	if err != nil {
		log.Fatal(err)
	}

	t := Trove{
		config:    c,
		producers: newProducers(c.Brokers),
		consumers: newConsumers(c.Brokers, c.GroupID),
		ext_match: matcher,
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
			ctx.SetBody([]byte(templates.PrintPage(&templates.Root{})))
		case "/dashboard":
			ctx.SetBody([]byte(templates.PrintPage(&templates.Dashboard{})))
		case "/ws":
			t.consumers.analytics.ws_handler(ctx)
		default:
			ctx.SetBody([]byte(templates.PrintPage(&templates.Unfound{})))
			err = errors.New("404")
		}
	}

	t.producers.logRequest(ctx, float64(time.Since(started))/float64(time.Second), err)
}

func (t *Trove) Run() {
	t.consumers.exitCallback = t.consumers.consume()

	var (
		ln  net.Listener
		err error
	)
	ln, err = net.Listen("tcp", t.config.Address)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Serving TLS connections on " + t.config.Address)
	err = t.server.ServeTLS(ln, t.config.Certificate, t.config.Key)

	if err != nil {
		log.Fatal(err)
	}
}

func (t *Trove) Exit() {
	if t.consumers.exitCallback != nil {
		t.consumers.exitCallback()
	}
	t.server.Shutdown()
}
