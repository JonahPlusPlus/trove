package main

import (
	"fmt"
	"time"

	"syscall/js"
)

var close = make(chan interface{})

func main() {
	defer graceful_shutdown()
	fmt.Println("WebAssembly Loaded")
	window := js.Global().Get("window")

	host := window.Get("location").Get("host").String()
	fmt.Println("Host at " + host)

	ws := js.Global().Get("WebSocket").New("wss://" + host + "/ws")

	ws.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Opened WebSocket")
		return nil
	}))

	ws.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer graceful_shutdown()
		for _, value := range args {
			fmt.Printf("Analytics: %s\n", value.Get("data").String())
			data := js.Global().Get("JSON").Call("parse", value.Get("data"))
			window.Set("trove_analytics", data)
			set_analytics(data)
			window.Call("trove_update")
		}
		return nil
	}))

	ws.Call("addEventListener", "close", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Closing WebSocket")

		close <- nil

		return nil
	}))

	interval := window.Call("setInterval", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer graceful_shutdown()
		fmt.Println("Updating Analytics")

		num_requests := window.Get("trove_analytics").Get("num_requests")

		client_msg := js.Global().Get("JSON").Call("stringify", num_requests)

		ws.Call("send", client_msg)

		return nil
	}), 15000)

	defer func() {
		window.Call("clearInterval", interval)
	}()

loop:
	for {
		select {
		case <-close:
			break loop
		default:
			time.Sleep(time.Second * 15)
		}
	}
}

func set_analytics(value js.Value) {
	defer graceful_shutdown()
	window := js.Global().Get("window")
	window.Set("trove_request_method_keys", js.Global().Get("Object").Call("keys", value.Get("request_method")))
	window.Set("trove_request_method_values", js.Global().Get("Object").Call("values", value.Get("request_method")))
	window.Set("trove_request_host_keys", js.Global().Get("Object").Call("keys", value.Get("request_host")))
	window.Set("trove_request_host_values", js.Global().Get("Object").Call("values", value.Get("request_host")))
	window.Set("trove_request_path_keys", js.Global().Get("Object").Call("keys", value.Get("request_path")))
	window.Set("trove_request_path_values", js.Global().Get("Object").Call("values", value.Get("request_path")))
}

func graceful_shutdown() {
	if r := recover(); r != nil {
		fmt.Println("Recovering from panic:", r)
		fmt.Println("Shutting down")
		close <- nil
	}
}
