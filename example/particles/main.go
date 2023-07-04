// Original JavaScript code by Elton Kamami:
// https://codepen.io/eltonkamami/pen/ECrKd
// Ported to Go by Frederik Zipp. Original copyright:
//
// Copyright (c) 2020 by Elton Kamami (https://codepen.io/eltonkamami/pen/ECrKd)
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

// Particles is a non-interactive canvas demo showing colorful animated
// particles.
//
// Usage:
//
//	particles [-http address]
//
// Flags:
//
//	-http  HTTP service address (e.g., '127.0.0.1:8080' or just ':8080').
//	       The default is ':8080'.
package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/fzipp/canvas"
)

const particlesNum = 500

var colors = []color.RGBA{
	{R: 0xf3, G: 0x5d, B: 0x4f, A: 0xff},
	{R: 0xf3, G: 0x68, B: 0x49, A: 0xff},
	{R: 0xc0, G: 0xd9, B: 0x88, A: 0xff},
	{R: 0x6d, G: 0xda, B: 0xf1, A: 0xff},
	{R: 0xf1, G: 0xe8, B: 0x5b, A: 0xff},
}

func main() {
	http := flag.String("http", ":8080", "HTTP service address (e.g., '127.0.0.1:8080' or just ':8080')")
	flag.Parse()

	fmt.Println("Listening on " + httpLink(*http))
	err := canvas.ListenAndServe(*http, run, &canvas.Options{
		Title:          "Particles",
		Width:          500,
		Height:         500,
		PageBackground: color.Black,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *canvas.Context) {
	demo := newParticlesDemo(
		float64(ctx.CanvasWidth()),
		float64(ctx.CanvasHeight()),
		particlesNum)

	ctx.SetGlobalCompositeOperation(canvas.OpLighter)
	ctx.SetLineWidth(0.5)

	for {
		select {
		case event := <-ctx.Events():
			if _, ok := event.(canvas.CloseEvent); ok {
				return
			}
		default:
			demo.draw(ctx)
			ctx.Flush()
			time.Sleep(5 * time.Millisecond)
		}
	}
}

type particlesDemo struct {
	w, h      float64
	particles []*particle
}

func newParticlesDemo(width, height float64, n int) *particlesDemo {
	particles := make([]*particle, n)
	for i := 0; i < len(particles); i++ {
		particles[i] = randomParticle(width, height)
	}
	return &particlesDemo{
		w: width, h: height,
		particles: particles,
	}
}

func (p *particlesDemo) draw(ctx *canvas.Context) {
	ctx.ClearRect(0, 0, p.w, p.h)
	for _, p1 := range p.particles {
		factor := 1.0

		for _, p2 := range p.particles {
			if p1.rgba == p2.rgba && p1.dist(p2) < 50 {
				ctx.SetStrokeStyle(p1.rgba)
				ctx.BeginPath()
				ctx.MoveTo(p1.x, p1.y)
				ctx.LineTo(p2.x, p2.y)
				ctx.Stroke()
				factor++
			}
		}

		ctx.SetFillStyle(p1.rgba)
		ctx.SetStrokeStyle(p1.rgba)

		ctx.BeginPath()
		ctx.Arc(p1.x, p1.y, p1.rad*factor, 0, math.Pi*2, true)
		ctx.Fill()
		ctx.ClosePath()

		ctx.BeginPath()
		ctx.Arc(p1.x, p1.y, (p1.rad+5)*factor, 0, math.Pi*2, true)
		ctx.Stroke()
		ctx.ClosePath()

		p1.x += p1.vx
		p1.y += p1.vy

		if p1.x > p.w {
			p1.x = 0
		}
		if p1.x < 0 {
			p1.x = p.w
		}
		if p1.y > p.h {
			p1.y = 0
		}
		if p1.y < 0 {
			p1.y = p.h
		}
	}
}

type particle struct {
	rgba   color.Color
	x, y   float64
	vx, vy float64
	rad    float64
}

func randomParticle(w, h float64) *particle {
	p := &particle{}
	p.x = math.Round(rand.Float64() * w)
	p.y = math.Round(rand.Float64() * h)
	p.rad = math.Round(rand.Float64()*1) + 1
	p.rgba = colors[int(math.Round(rand.Float64()*3))]
	p.vx = math.Round(rand.Float64()*3) - 1.5
	p.vy = math.Round(rand.Float64()*3) - 1.5
	return p
}

func (p *particle) dist(other *particle) float64 {
	return math.Sqrt(math.Pow(other.x-p.x, 2) + math.Pow(other.y-p.y, 2))
}

func httpLink(addr string) string {
	if addr[0] == ':' {
		addr = "localhost" + addr
	}
	return "http://" + addr
}
