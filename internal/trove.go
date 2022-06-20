package trove

import (
	"io/ioutil"
	"log"
	"mime"
	"regexp"
	"strings"
	"time"

	"github.com/JonahPlusPlus/trove/templates"
	"github.com/valyala/fasthttp"
)

type Trove struct {
	config    Config
	producers Producers
	ext_match *regexp.Regexp
}

func New(config ...Config) Trove {
	c := getConfig(config...)

	matcher, err := regexp.Compile(`\..*`)

	if err != nil {
		log.Fatal(err)
	}

	return Trove{
		config:    c,
		producers: newProducers(c.Broker),
		ext_match: matcher,
	}
}

func (t *Trove) handler(ctx *fasthttp.RequestCtx) {
	started := time.Now()
	path := string(ctx.Path())

	if strings.Contains(path, ".") {
		// Handle static files
		f, err := ioutil.ReadFile("./static" + path)
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
		default:
			ctx.SetBody([]byte(templates.PrintPage(&templates.Unfound{})))
		}
	}

	t.producers.logRequest(ctx, float64(time.Since(started))/float64(time.Second))
}

func (t *Trove) Run() {
	log.Fatal(fasthttp.ListenAndServeTLS(t.config.Address, t.config.Certificate, t.config.Key, t.handler))
}
