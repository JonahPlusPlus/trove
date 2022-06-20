package trove

import (
	"strings"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Trove struct {
	config    Config
	producers Producers
	router    *router.Router
}

func New(config ...Config) Trove {
	c := getConfig(config...)
	router := router.New()

	router.ServeFiles("/{filepath:*}", "./static/")

	return Trove{
		config:    c,
		producers: newProducers(c.Broker),
		router:    router,
	}
}

func (t *Trove) Address() string {
	return t.config.Address
}

func (t *Trove) CertificatePath() string {
	return t.config.Certificate
}

func (t *Trove) KeyPath() string {
	return t.config.Key
}

func (t *Trove) Run(ctx *fasthttp.RequestCtx) {
	started := time.Now()
	path := string(ctx.Path())

	if !strings.Contains(path, ".") {
		var index_path string

		// check for trailing slash
		if path[len(path)-1] == '/' {
			index_path = "index.html"
		} else {
			index_path = "/index.html"
		}

		ctx.URI().SetPath(path + index_path)
	}

	t.router.Handler(ctx)

	t.producers.logRequest(ctx, float64(time.Since(started))/float64(time.Second))

}
