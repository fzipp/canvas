// Copyright 2021 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example was ported from:
// https://codepen.io/hakimel/pen/KanIi
// Original copyright:
//
// Copyright (c) 2021 by Hakim El Hattab (https://codepen.io/hakimel/pen/KanIi)
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

// Trail is an interactive animation of rotating particles that follow the
// mouse or touch pointer.
//
// Usage:
//
//	trail [-http address]
//
// Flags:
//
//	-http  HTTP service address (e.g., '127.0.0.1:8080' or just ':8080').
//	       The default is ':8080'.
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"time"

	"github.com/fzipp/canvas"
)

const (
	radius = 110

	radiusScaleMin = 1
	radiusScaleMax = 1.5

	// The number of particles that are used to generate the trail
	quantity = 25
)

func main() {
	http := flag.String("http", ":8080", "HTTP service address (e.g., '127.0.0.1:8080' or just ':8080')")
	flag.Parse()

	fmt.Println("Listening on " + httpLink(*http))
	err := canvas.ListenAndServe(*http, run, &canvas.Options{
		Width:             750,
		Height:            1334,
		ScaleToPageWidth:  true,
		ScaleToPageHeight: true,
		EnabledEvents: []canvas.Event{
			canvas.MouseMoveEvent{},
			canvas.MouseDownEvent{},
			canvas.MouseUpEvent{},
			canvas.TouchStartEvent{},
			canvas.TouchMoveEvent{},
		},
		ReconnectInterval: time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *canvas.Context) {
	d := &demo{
		mouseX:      float64(ctx.CanvasWidth() / 2),
		mouseY:      float64(ctx.CanvasHeight() / 2),
		radiusScale: 1.0,
	}
	d.createParticles()
	for {
		select {
		case event := <-ctx.Events():
			if _, ok := event.(canvas.CloseEvent); ok {
				return
			}
			d.handle(event)
		default:
			d.draw(ctx)
			ctx.Flush()
			time.Sleep((1000 / 60) * time.Millisecond)
		}
	}
}

type point struct {
	x, y float64
}

type particle struct {
	position   point
	shift      point
	size       float64
	angle      float64
	speed      float64
	targetSize float64
	fillColor  string
	orbit      float64
}

type demo struct {
	particles   []particle
	radiusScale float64
	mouseX      float64
	mouseY      float64
	mouseIsDown bool
}

func (d *demo) createParticles() {
	d.particles = make([]particle, 0, quantity)
	for i := 0; i < quantity; i++ {
		p := particle{
			position:   point{x: d.mouseX, y: d.mouseY},
			shift:      point{x: d.mouseX, y: d.mouseY},
			size:       1,
			angle:      0,
			speed:      0.01 + rand.Float64()*0.04,
			targetSize: 1,
			fillColor:  "#" + fmt.Sprintf("%6x", int(rand.Float64()*0x404040+0xaaaaaa)),
			orbit:      radius*.5 + (radius * .5 * rand.Float64()),
		}
		d.particles = append(d.particles, p)
	}
}

func (d *demo) handle(ev canvas.Event) {
	switch e := ev.(type) {
	case canvas.MouseMoveEvent:
		d.mouseX = float64(e.X)
		d.mouseY = float64(e.Y)
	case canvas.MouseDownEvent:
		d.mouseIsDown = true
	case canvas.MouseUpEvent:
		d.mouseIsDown = false
	case canvas.TouchStartEvent:
		if len(e.Touches) == 1 {
			d.mouseX = float64(e.Touches[0].X)
			d.mouseY = float64(e.Touches[0].Y)
		}
	case canvas.TouchMoveEvent:
		if len(e.Touches) == 1 {
			d.mouseX = float64(e.Touches[0].X)
			d.mouseY = float64(e.Touches[0].Y)
		}
	}
}

func (d *demo) draw(ctx *canvas.Context) {
	if d.mouseIsDown {
		// Scale upward to the max scale
		d.radiusScale += (radiusScaleMax - d.radiusScale) * (0.02)
	} else {
		// Scale downward to the min scale
		d.radiusScale -= (d.radiusScale - radiusScaleMin) * (0.02)
	}

	d.radiusScale = math.Min(d.radiusScale, radiusScaleMax)

	// Fade out the lines slowly by drawing a rectangle over the entire canvas
	ctx.SetFillStyleString("rgba(0,0,0,0.05)")
	ctx.FillRect(0, 0, float64(ctx.CanvasWidth()), float64(ctx.CanvasHeight()))

	for i := range d.particles {
		p := &d.particles[i]

		lp := point{x: p.position.x, y: p.position.y}

		// Offset the angle to keep the spin going
		p.angle += p.speed

		// Follow mouse with some lag
		p.shift.x += (d.mouseX - p.shift.x) * (p.speed)
		p.shift.y += (d.mouseY - p.shift.y) * (p.speed)

		// Apply position
		p.position.x = p.shift.x + math.Cos(float64(i)+p.angle)*(p.orbit*d.radiusScale)
		p.position.y = p.shift.y + math.Sin(float64(i)+p.angle)*(p.orbit*d.radiusScale)

		// Limit to screen bounds
		p.position.x = math.Max(math.Min(p.position.x, float64(ctx.CanvasWidth())), 0)
		p.position.y = math.Max(math.Min(p.position.y, float64(ctx.CanvasHeight())), 0)

		p.size += (p.targetSize - p.size) * 0.05

		// If we're at the target size, set a new one. Think of it like a regular day at work.
		if math.Round(p.size) == math.Round(p.targetSize) {
			p.targetSize = 1 + rand.Float64()*7
		}

		ctx.BeginPath()
		ctx.SetFillStyleString(p.fillColor)
		ctx.SetStrokeStyleString(p.fillColor)
		ctx.SetLineWidth(p.size)
		ctx.MoveTo(lp.x, lp.y)
		ctx.LineTo(p.position.x, p.position.y)
		ctx.Stroke()
		ctx.Arc(p.position.x, p.position.y, p.size/2, 0, math.Pi*2, true)
		ctx.Fill()
	}
}

func httpLink(addr string) string {
	if addr[0] == ':' {
		addr = "localhost" + addr
	}
	return "http://" + addr
}
