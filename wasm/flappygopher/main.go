//go:build !solution

package main

import (
	"syscall/js"
	"time"
)

var (
	canvas  js.Value
	context js.Value
	width   int
	height  int

	gopherImage js.Value
	gopherX     float64
	gopherY     float64
	velocity    float64
	gravity     float64 = 0.6
	lift        float64 = -15
)

func main() {
	c := make(chan struct{}, 0)

	doc := js.Global().Get("document")
	canvas = doc.Call("getElementById", "gameCanvas")
	if canvas.IsNull() {
		println("Canvas element not found")
		close(c)
		return
	}

	context = canvas.Call("getContext", "2d")
	width = canvas.Get("width").Int()
	height = canvas.Get("height").Int()

	gopherImage = js.Global().Get("Image").New()
	gopherImage.Set("src", "gopher.png")

	js.Global().Get("window").Call("addEventListener", "click", js.FuncOf(flyUp))

	gopherImage.Call("addEventListener", "load", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go gameLoop()
		return nil
	}))

	<-c
}

func drawGopher() {
	context.Call("drawImage", gopherImage, gopherX, gopherY)
}

func flyUp(this js.Value, p []js.Value) interface{} {
	velocity += lift
	return nil
}

func gameLoop() {
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		update()
	}
}

func update() {
	velocity += gravity
	velocity *= 0.9 // Фактор трения
	gopherY += velocity

	if gopherY+float64(gopherImage.Get("height").Int()) > float64(height) {
		gopherY = float64(height) - float64(gopherImage.Get("height").Int())
		velocity = 0
	}

	if gopherY < 0 {
		gopherY = 0
		velocity = 0
	}

	context.Call("clearRect", 0, 0, width, height)
	drawGopher()
}
