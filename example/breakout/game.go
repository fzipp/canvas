// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/fzipp/canvas"
)

type game struct {
	started bool
	quit    bool
	score   int
	round   int
	size    vec2
	bricks  []brick
	paddle  paddle
	ball    ball
}

func newGame(size vec2) *game {
	g := &game{size: size}
	g.resetGame()
	return g
}

func (g *game) resetGame() {
	g.started = true
	g.score = 0
	g.round = 1
	paddleSize := vec2{x: 100, y: 20}
	g.paddle = paddle{
		pos:   vec2{x: g.size.x / 2, y: g.size.y - (paddleSize.y / 2)},
		size:  paddleSize,
		color: color.RGBA{R: 0x0a, G: 0x85, B: 0xc2, A: 0xff},
	}
	ballRadius := 5.0
	g.ball = ball{
		radius: ballRadius,
		v:      vec2{x: 1, y: -1},
		color:  color.White,
	}
	g.bricks = g.initialBricks(14, 8)
	g.resetBall()
}

var rowGroupColors = []color.Color{
	color.RGBA{R: 0xa3, G: 0x1e, B: 0x0a, A: 0xff},
	color.RGBA{R: 0xc2, G: 0x85, B: 0x0a, A: 0xff},
	color.RGBA{R: 0x0a, G: 0x85, B: 0x33, A: 0xff},
	color.RGBA{R: 0xc2, G: 0xc2, B: 0x29, A: 0xff},
}

var rowGroupPoints = []int{7, 5, 3, 1}

func (g *game) initialBricks(columns, rows int) []brick {
	width := int(g.size.x) / columns
	height := 30
	bricks := make([]brick, 0, columns*rows)
	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			x := col * width
			y := row * height
			rowGroup := (row / 2) % len(rowGroupColors)
			bricks = append(bricks, brick{
				rect:   image.Rect(x+1, y+1, x+width-1, y+height-1),
				color:  rowGroupColors[rowGroup],
				points: rowGroupPoints[rowGroup],
			})
		}
	}
	return bricks
}

func (g *game) handle(event canvas.Event) {
	switch e := event.(type) {
	case canvas.CloseEvent:
		g.quit = true
	case canvas.MouseMoveEvent:
		if g.started {
			g.paddle.pos.x = float64(e.X)
		}
	case canvas.TouchStartEvent:
		if len(e.Touches) == 1 {
			g.paddle.pos.x = float64(e.Touches[0].X)
		}
	case canvas.TouchMoveEvent:
		if len(e.Touches) == 1 {
			g.paddle.pos.x = float64(e.Touches[0].X)
		}
	case canvas.KeyDownEvent:
		const paddleSpeedX = 15
		switch e.Key {
		case "ArrowRight":
			g.paddle.pos.x += paddleSpeedX
			if g.paddle.pos.x >= g.size.x {
				g.paddle.pos.x = g.size.x - 1
			}
		case "ArrowLeft":
			g.paddle.pos.x -= paddleSpeedX
			if g.paddle.pos.x < 0 {
				g.paddle.pos.x = 0
			}
		case " ":
			g.started = !g.started
		}
	}
}

func (g *game) update() {
	if !g.started {
		return
	}
	g.ball.update()
	g.checkWallCollisions()
	g.checkBrickCollisions()
	g.checkPaddleCollision()
}

func (g *game) resetBall() {
	g.ball.pos = g.paddle.pos.sub(vec2{x: 0, y: g.ball.radius + (g.paddle.size.y / 2)})
}

func (g *game) checkWallCollisions() {
	ballBounds := g.ball.bounds()
	gameBounds := g.bounds()
	bottom := ballBounds.Max.Y >= gameBounds.Max.Y
	left := ballBounds.Min.X <= gameBounds.Min.X
	right := ballBounds.Max.X >= gameBounds.Max.X
	top := ballBounds.Min.Y <= gameBounds.Min.Y
	if left || right {
		g.ball.v.x = -g.ball.v.x
	} else if top {
		g.ball.v.y = -g.ball.v.y
	} else if bottom {
		g.round++
		if g.round > 3 {
			g.resetGame()
		} else {
			g.resetBall()
		}
	}
}

func (g *game) checkBrickCollisions() {
	survivingBricks := make([]brick, 0, len(g.bricks))
	for _, brick := range g.bricks {
		collision := g.ball.bounceOnCollision(brick.bounds())
		if collision == collisionNone {
			survivingBricks = append(survivingBricks, brick)
		} else {
			g.score += brick.points
		}
	}
	g.bricks = survivingBricks
}

func (g *game) checkPaddleCollision() collision {
	return g.ball.bounceOnCollision(g.paddle.bounds())
}

func (g *game) draw(ctx *canvas.Context) {
	g.drawBackground(ctx)
	for _, brick := range g.bricks {
		brick.draw(ctx)
	}
	g.paddle.draw(ctx)
	g.ball.draw(ctx)
	g.drawScore(ctx)
}

func (g *game) drawBackground(ctx *canvas.Context) {
	ctx.SetFillStyle(color.Black)
	ctx.FillRect(0, 0, g.size.x, g.size.y)
}

func (g *game) drawScore(ctx *canvas.Context) {
	ctx.SetFillStyle(color.White)
	ctx.FillText(fmt.Sprintf("%03d    %d", g.score, g.round), 10, 35)
}

func (g *game) bounds() image.Rectangle {
	return image.Rect(0, 0, int(g.size.x), int(g.size.y))
}
