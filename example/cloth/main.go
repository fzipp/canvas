// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example was ported from:
// https://codepen.io/dissimulate/pen/KrAwx
// Original copyright:
//
// Copyright (c) 2020 by dissimulate (https://codepen.io/dissimulate/pen/KrAwx)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/fzipp/canvas"
)

const (
	physicsAccuracy = 3
	mouseInfluence  = 20
	mouseCut        = 5
	gravity         = 1200
	clothHeight     = 30
	clothWidth      = 50
	startY          = 20
	spacing         = 7
	tearDistance    = 60
)

func main() {
	port := ":8080"
	fmt.Println("Listening on http://localhost" + port)
	canvas.ListenAndServe(port, run,
		canvas.Size(560, 350),
		canvas.Title("Tearable Cloth"),
		canvas.EnableEvents(
			canvas.SendMouseMove,
			canvas.SendMouseDown,
			canvas.SendMouseUp,
		),
	)
}

func run(ctx *canvas.Context) {
	ctx.SetStrokeStyle(color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF})

	cloth := newCloth(
		ctx.CanvasWidth(),
		float64(ctx.CanvasWidth()-1),
		float64(ctx.CanvasHeight()-1),
	)

	for {
		select {
		case <-ctx.Quit():
			return
		case event := <-ctx.Events():
			cloth.handle(event)
		default:
			cloth.update()
			cloth.draw(ctx)
			ctx.Flush()
			time.Sleep(5 * time.Millisecond)
		}
	}
}

type Cloth struct {
	boundsX float64
	boundsY float64
	mouse   Mouse
	points  []*Point
}

func newCloth(canvasWidth int, boundsX, boundsY float64) *Cloth {
	cloth := &Cloth{
		boundsX: boundsX,
		boundsY: boundsY,
	}
	startX := float64(canvasWidth)/2 - clothWidth*spacing/2
	for y := 0; y <= clothHeight; y++ {
		for x := 0; x <= clothWidth; x++ {
			p := newPoint(
				startX+float64(x*spacing),
				startY+float64(y*spacing),
			)
			if x != 0 {
				p.attach(cloth.points[len(cloth.points)-1])
			}
			if y == 0 {
				p.pin(p.x, p.y)
			}
			if y != 0 {
				p.attach(cloth.points[x+(y-1)*(clothWidth+1)])
			}
			cloth.points = append(cloth.points, p)
		}
	}
	return cloth
}

func (c *Cloth) handle(event canvas.Event) {
	switch e := event.(type) {
	case canvas.MouseMoveEvent:
		c.mouse.px = c.mouse.x
		c.mouse.py = c.mouse.y
		c.mouse.x = float64(e.X)
		c.mouse.y = float64(e.Y)
	case canvas.MouseUpEvent:
		c.mouse.down = false
	case canvas.MouseDownEvent:
		c.mouse.button = e.Button
		c.mouse.px = c.mouse.x
		c.mouse.py = c.mouse.y
		c.mouse.x = float64(e.X)
		c.mouse.y = float64(e.Y)
		c.mouse.down = true
	}
}

func (c *Cloth) update() {
	for i := 0; i < physicsAccuracy; i++ {
		for _, p := range c.points {
			p.resolveConstraints(c.boundsX, c.boundsY)
		}
	}
	for _, p := range c.points {
		p.update(.016, &c.mouse)
	}
}

func (c *Cloth) draw(ctx *canvas.Context) {
	ctx.ClearRect(0, 0,
		float64(ctx.CanvasWidth()),
		float64(ctx.CanvasHeight()))

	ctx.BeginPath()
	for _, p := range c.points {
		p.draw(ctx)
	}
	ctx.Stroke()
}

type Mouse struct {
	down   bool
	button int
	x, y   float64
	px, py float64
}

type Point struct {
	x, y        float64
	px, py      float64
	vx, vy      float64
	pinX, pinY  float64
	constraints []*Constraint
}

func newPoint(x, y float64) *Point {
	return &Point{
		x: x, y: y,
		px: x, py: y,
		vx: 0, vy: 0,
		pinX: math.NaN(), pinY: math.NaN(),
	}
}

func (p *Point) update(delta float64, mouse *Mouse) {
	if mouse.down {
		diffX := p.x - mouse.x
		diffY := p.y - mouse.y
		dist := math.Sqrt(diffX*diffX + diffY*diffY)

		if mouse.button == 1 {
			if dist < mouseInfluence {
				p.px = p.x - (mouse.x-mouse.px)*1.8
				p.py = p.y - (mouse.y-mouse.py)*1.8
			}
		} else if dist < mouseCut {
			p.constraints = p.constraints[:0]
		}
	}

	p.addForce(0, gravity)

	delta *= delta
	nx := p.x + ((p.x - p.px) * .99) + ((p.vx / 2) * delta)
	ny := p.y + ((p.y - p.py) * .99) + ((p.vy / 2) * delta)

	p.px = p.x
	p.py = p.y

	p.x = nx
	p.y = ny

	p.vy = 0
	p.vx = 0
}

func (p *Point) draw(ctx *canvas.Context) {
	for _, c := range p.constraints {
		c.draw(ctx)
	}
}

func (p *Point) attach(q *Point) {
	p.constraints = append(p.constraints, newConstraint(p, q))
}

func (p *Point) pin(x, y float64) {
	p.pinX = x
	p.pinY = y
}

func (p *Point) addForce(x, y float64) {
	p.vx += x
	p.vy += y

	const round = 400
	p.vx = math.Floor(p.vx*round) / round
	p.vy = math.Floor(p.vy*round) / round
}

func (p *Point) resolveConstraints(boundsX, boundsY float64) {
	if !math.IsNaN(p.pinX) && !math.IsNaN(p.pinY) {
		p.x = p.pinX
		p.y = p.pinY
		return
	}

	for _, c := range p.constraints {
		c.resolve()
	}

	if p.x > boundsX {
		p.x = 2*boundsX - p.x
	} else {
		if 1 > p.x {
			p.x = 2 - p.x
		}
	}
	if p.y < 1 {
		p.y = 2 - p.y
	} else {
		if p.y > boundsY {
			p.y = 2*boundsY - p.y
		}
	}
}

func (p *Point) removeConstraint(c *Constraint) {
	for i, elem := range p.constraints {
		if elem == c {
			p.constraints = append(p.constraints[:i], p.constraints[i+1:]...)
			return
		}
	}
}

type Constraint struct {
	p1, p2 *Point
	length float64
}

func newConstraint(p1, p2 *Point) *Constraint {
	return &Constraint{p1: p1, p2: p2, length: spacing}
}

func (c *Constraint) draw(ctx *canvas.Context) {
	ctx.MoveTo(c.p1.x, c.p1.y)
	ctx.LineTo(c.p2.x, c.p2.y)
}

func (c *Constraint) resolve() {
	diffX := c.p1.x - c.p2.x
	diffY := c.p1.y - c.p2.y
	dist := math.Sqrt(diffX*diffX + diffY*diffY)
	diff := (c.length - dist) / dist

	if dist > tearDistance {
		c.p1.removeConstraint(c)
	}

	px := diffX * diff * 0.5
	py := diffY * diff * 0.5

	c.p1.x += px
	c.p1.y += py
	c.p2.x -= px
	c.p2.y -= py
}
