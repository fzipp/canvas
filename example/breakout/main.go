// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/fzipp/canvas"
)

func main() {
	port := ":8080"
	fmt.Println("Listening on http://localhost" + port)
	err := canvas.ListenAndServe(port, run,
		canvas.Size(1334, 750),
		canvas.FullPage(),
		canvas.Title("Breakout"),
		canvas.DisableCursor(),
		canvas.EnableEvents(
			canvas.MouseMoveEvent{},
			canvas.KeyPressEvent{},
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
	size := image.Pt(ctx.CanvasWidth(), ctx.CanvasHeight())
	game := newGame(size)
	ctx.SetFont("30px sans-serif")
	for {
		select {
		case <-ctx.Quit():
			return
		case event := <-ctx.Events():
			game.handle(event, ctx)
		default:
			game.update()
			game.draw(ctx)
			ctx.Flush()
			time.Sleep(5 * time.Millisecond)
		}
	}
}
