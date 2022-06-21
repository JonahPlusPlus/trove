package trove

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

type Analytics struct {
	// Number of requests made
	NumRequests *uint `json:"num_requests"`
	// Pie chart of HTTP Method
	RequestMethod map[string]uint `json:"request_method"`
	// Pie chart of Host
	RequestHost map[string]uint `json:"request_host"`
	// Pie chart of Paths
	RequestPath map[string]uint `json:"request_path"`
	// Bar chart of time
	RequestTime *float64 `json:"request_time"`
}

func newAnalytics() Analytics {
	var num_requests uint = 0
	var request_time float64 = 0
	return Analytics{
		NumRequests:   &num_requests,
		RequestMethod: make(map[string]uint),
		RequestHost:   make(map[string]uint),
		RequestPath:   make(map[string]uint),
		RequestTime:   &request_time,
	}
}

func (a *Analytics) addRequestEvent(event RequestEvent) {
	new_avg := *a.RequestTime * float64(*a.NumRequests)

	*a.NumRequests += 1
	a.RequestMethod[event.Method]++
	a.RequestHost[event.Host]++
	a.RequestPath[event.Path]++
	*a.RequestTime = (new_avg + event.Time) / float64(*a.NumRequests)
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 30 * time.Second
	pingPeriod     = pongWait / 2
	maxMessageSize = 512
)

var upgrader = websocket.FastHTTPUpgrader{} // use default options

func (a *Analytics) ws_handler(ctx *fasthttp.RequestCtx) {
	pingTicker := time.NewTicker(pingPeriod)
	exit := make(chan interface{})
	err := upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
		defer func() {
			pingTicker.Stop()
			ws.Close()
		}()
		ws.WriteJSON(a)
		go func() {
			defer func() {
				exit <- nil
			}()
			for {
				_, msg, err := ws.ReadMessage()

				if err != nil {
					log.Println("WS reading error:", err)
					break
				}

				i, err := strconv.ParseInt(string(msg), 10, 0)

				if err != nil {
					log.Println("WS parse error:", err)
					break
				}

				if uint(i) == *a.NumRequests {
					continue
				}

				ws.SetWriteDeadline(time.Now().Add(writeWait))
				b, err := json.Marshal(a)
				if err != nil {
					log.Println("WS marshal error:", err)
					break
				}
				if err := ws.WriteMessage(websocket.TextMessage, b); err != nil {
					log.Println("WS writing error:", err)
					break
				}
			}
		}()
	ping:
		for {
			select {
			case <-pingTicker.C:
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					log.Println("WS error:", err)
					return
				}
			case <-exit:
				break ping
			}
		}

	})

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			log.Println(err)
		}
		return
	}
}
