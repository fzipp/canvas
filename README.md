# canvas

[![PkgGoDev](https://pkg.go.dev/badge/github.com/fzipp/canvas)](https://pkg.go.dev/github.com/fzipp/canvas)
![Build Status](https://github.com/fzipp/canvas/workflows/build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fzipp/canvas)](https://goreportcard.com/report/github.com/fzipp/canvas)

This Go module utilizes WebSockets to establish communication with a
[2D canvas graphics context](https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D)
in a web browser,
providing a portable way to create interactive 2D graphics
from within a Go program.

The Go program (server) sends draw commands to the web browser (client)
via WebSocket using a binary format.
In return, the client sends keyboard, mouse, and touch events to the server.

This module does not rely on operating system-specific backends
or Cgo bindings.
It also does not utilize WebAssembly,
which means the Go code runs on the server side,
rather than in the browser.
The client-server design enables the canvas
to be displayed on a different machine over the network.

## Examples

The [example](example) subdirectory contains a variety of demo programs.

![Screenshots of examples](https://github.com/fzipp/canvas/blob/assets/examples.png)

## Usage

### Drawing

The `ListenAndServe` function initializes the canvas server
and takes the following arguments:
the network address with the port number to bind to,
a run function,
and an options structure that configures various aspects
such as the canvas size in pixels
or a title for the browser tab.

The `run` function is called when a client connects to the server.
This serves as the entry point for drawing.

```go
package main

import (
	"image/color"
	"log"

	"github.com/fzipp/canvas"
)

func main() {
	err := canvas.ListenAndServe(":8080", run, &canvas.Options{
		Title:  "Example 1: Drawing",
		Width:  100,
		Height: 80,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *canvas.Context) {
	ctx.SetFillStyle(color.RGBA{R: 200, A: 255})
	ctx.FillRect(10, 10, 50, 50)
	// ...
	ctx.Flush()
}
```

After starting the program,
you can access the canvas by opening http://localhost:8080
in a web browser.

The server doesn't immediately send each drawing operation to the client
but instead buffers them until the `Flush` method is called.
The flush should occur once the image or an animation frame is complete;
otherwise, nothing will be displayed.

Each client connection starts its own run function as a goroutine.
Access to shared state between client connections must be synchronized.
If you don't want to share state between connections,
you should keep it local to the run function
and pass the state to other functions called by the run function.

### An animation loop

To create an animation,
you can use a `for` loop within the `run` function.
Inside this loop,
observe the `ctx.Events()` channel
for a `canvas.CloseEvent` to exit the loop
when the connection is closed.

A useful pattern is to create a struct
that holds the animation state
and has both an update and a draw method:

```go
package main

import (
	"log"
	"time"

	"github.com/fzipp/canvas"
)

func main() {
	err := canvas.ListenAndServe(":8080", run, &canvas.Options{
		Title:  "Example 2: Animation",
		Width:  800,
		Height: 600,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *canvas.Context) {
	d := &demo{}
	for {
		select {
		case event := <-ctx.Events():
			if _, ok := event.(canvas.CloseEvent); ok {
				return
			}
		default:
			d.update()
			d.draw(ctx)
			ctx.Flush()
			time.Sleep(time.Second / 6)
		}
	}
}

type demo struct {
	// Animation state, for example:
	x, y int
	// ...
}

func (d *demo) update() {
	// Update animation state for the next frame
	// ...
}

func (d *demo) draw(ctx *canvas.Context) {
	// Draw the frame here, based on the animation state
	// ...
}
```

### Keyboard, mouse and touch events

To handle keyboard, mouse, and touch events,
you need to specify which events the client should observe
and send to the server.
This is achieved by passing an `EnabledEvents` option
to the `ListenAndServe` function.
Mouse move events typically generate more WebSocket communication
than the others,
so you may want to enable them only if necessary.

The `ctx.Events()` channel receives the observed events,
and a type switch is used to determine the specific event type.
A useful pattern involves creating a `handle` method for event handling:

```go
package main

import (
	"log"

	"github.com/fzipp/canvas"
)

func main() {
	err := canvas.ListenAndServe(":8080", run, &canvas.Options{
		Title:  "Example 3: Events",
		Width:  800,
		Height: 600,
		EnabledEvents: []canvas.Event{
			canvas.MouseDownEvent{},
			canvas.MouseMoveEvent{},
			canvas.TouchStartEvent{},
			canvas.TouchMoveEvent{},
			canvas.KeyDownEvent{},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *canvas.Context) {
	d := &demo{}
	for !d.quit {
		select {
		case event := <-ctx.Events():
			d.handle(event)
		default:
			d.update()
			d.draw(ctx)
			ctx.Flush()
		}
	}
}

type demo struct {
	quit bool
	// ...
}

func (d *demo) handle(event canvas.Event) {
	switch e := event.(type) {
	case canvas.CloseEvent:
		d.quit = true
	case canvas.MouseDownEvent:
		// ...
	case canvas.MouseMoveEvent:
		// ...
	case canvas.TouchStartEvent:
		// ...
   	case canvas.TouchMoveEvent:
		// ...
   	case canvas.KeyDownEvent:
		// ...
	}
}

func (d *demo) update() {
	// ...
}

func (d *demo) draw(ctx *canvas.Context) {
	// ...
}
```

Note that the `canvas.CloseEvent` does not have to be explicitly enabled.
It is always enabled by default.

## Alternatives

* [github.com/tfriedel6/canvas](https://github.com/tfriedel6/canvas) -
  A canvas implementation for Go with OpenGL backends for various
  operating systems.
* [github.com/llgcode/draw2d](https://github.com/llgcode/draw2d) -
  A 2D vector graphics library for Go with support for multiple outputs
  such as images, PDF documents, OpenGL and SVG.
* [github.com/ajstarks/svgo](https://github.com/ajstarks/svgo) -
  A Go library for SVG generation.
* [github.com/tdewolff/canvas](https://github.com/tdewolff/canvas) -
  A common vector drawing target that can output SVG, PDF, EPS,
  raster images (PNG, JPG, GIF, ...), HTML Canvas through WASM, and OpenGL.
* [github.com/fogleman/gg](https://github.com/fogleman/gg) -
  A library for rendering 2D graphics in pure Go.

2D game engines:

* [github.com/faiface/pixel](https://github.com/faiface/pixel) - Pixel
* [github.com/hajimehoshi/ebiten](https://github.com/hajimehoshi/ebiten) - Ebiten
* [github.com/oakmound/oak](https://github.com/oakmound/oak) - Oak

## License

This project is free and open source software licensed under the
[BSD 3-Clause License](LICENSE).
