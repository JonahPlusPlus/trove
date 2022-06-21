package main

import (
	"fmt"

	"syscall/js"
)

func main() {
	fmt.Println("WebAssembly Loaded")
	window := js.Global().Get("window")
	fmt.Println("Host at " + window.Get("location").Get("host").String())
	ws := js.Global().Get("WebSocket").New("wss://" + js.Global().Get("window").Get("location").Get("host").String() + "/ws")

	ws.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Opened WebSocket")
		return nil
	}))

	ws.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		for _, value := range args {
			fmt.Printf("Analytics: %s\n", value.Get("data").String())
			js.Global().Get("window").Set("trove_analytics", js.Global().Get("JSON").Call("parse", value.Get("data")))
		}
		return nil
	}))

	close := make(chan interface{})

	ws.Call("addEventListener", "close", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Closing WebSocket")

		close <- nil

		return nil
	}))

	window.Call("setInterval", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Updating Analytics")

		num_requests := window.Get("trove_analytics").Get("num_requests")

		client_msg := js.Global().Get("JSON").Call("stringify", num_requests)

		ws.Call("send", client_msg)

		return nil
	}), 15000)

	<-close

}
