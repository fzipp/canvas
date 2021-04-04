// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A simple Breakout game.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/fzipp/canvas"
)

func main() {
	http := flag.String("http", ":8080", "HTTP service address (e.g., '127.0.0.1:8080' or just ':8080')")
	flag.Parse()

	fmt.Println("Listening on " + httpLink(*http))
	err := canvas.ListenAndServe(*http, run,
		canvas.Size(1334, 750),
		canvas.ScaleFullPage(true, true),
		canvas.Title("Breakout"),
		canvas.DisableCursor(),
		canvas.EnableEvents(
			canvas.MouseMoveEvent{},
			canvas.KeyDownEvent{},
			canvas.TouchStartEvent{},
			canvas.TouchMoveEvent{},
		),
		canvas.Reconnect(time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *canvas.Context) {
	size := vec2{x: float64(ctx.CanvasWidth()), y: float64(ctx.CanvasHeight())}
	game := newGame(size)
	ctx.SetFont("30px sans-serif")
	for !game.quit {
		select {
		case event := <-ctx.Events():
			game.handle(event)
		default:
			game.update()
			game.draw(ctx)
			ctx.Flush()
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func httpLink(addr string) string {
	if addr[0] == ':' {
		addr = "localhost" + addr
	}
	return "http://" + addr
}
