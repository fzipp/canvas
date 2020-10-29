// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"time"

	"github.com/fzipp/canvas"
)

func main() {
	canvas.ListenAndServe(":8080",
		800, 600, "Breakout", run,
		canvas.SendMouseMove|canvas.SendKeyPress|canvas.SendKeyDown)
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
			game.handle(event)
		default:
			game.update()
			game.draw(ctx)
			ctx.Flush()
			time.Sleep(5 * time.Millisecond)
		}
	}
}
