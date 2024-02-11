// Original JavaScript code by Victor Ribeiro:
// https://github.com/victorqribeiro/retroSynthwave
// Ported to Go by Frederik Zipp. Original copyright:
//
// MIT License
//
// Copyright (c) 2021 Victor Ribeiro
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Retro-synthwave is an animated demo by Victor Ribeiro.
// It shows a flight over a grid-like mountain landscape with a sunset backdrop.
//
// Usage:
//
//	retro-synthwave [-http address]
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
	"math/rand/v2"
	"time"

	"github.com/fzipp/canvas"
)

func main() {
	http := flag.String("http", ":8080", "HTTP service address (e.g., '127.0.0.1:8080' or just ':8080')")
	flag.Parse()

	fmt.Println("Listening on " + httpLink(*http))
	err := canvas.ListenAndServe(*http, run, &canvas.Options{
		Title:             "Retro Synthwave",
		Width:             1440,
		Height:            694,
		ScaleToPageHeight: true,
		PageBackground:    color.Black,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type demo struct {
	w, h       float64
	points     [][]point
	offset     float64
	spacing    float64
	gradient   *canvas.Gradient
	background *canvas.Gradient
}

type point struct {
	x, y, z float64
}

func run(ctx *canvas.Context) {
	d := &demo{
		w: float64(ctx.CanvasWidth()),
		h: float64(ctx.CanvasHeight()),
	}
	d.spacing = 40.0
	d.points = make([][]point, 30)
	for i := range d.points {
		d.points[i] = make([]point, 60)
		for j := range d.points[i] {
			dist := math.Abs(float64(j) - float64(len(d.points[0]))/2)
			d.points[i][j] = point{
				x: float64(j) * d.spacing,
				y: rand.Float64()*-(dist*dist) + 30,
				z: -float64(i) * 10,
			}
		}
	}
	d.offset = float64(len(d.points[0])) * d.spacing / 2

	d.gradient = ctx.CreateLinearGradient(0, -150, 0, 100)
	d.gradient.AddColorStopString(0, "gold")
	d.gradient.AddColorStopString(1, "rgb(200, 0, 100)")
	defer d.gradient.Release()

	d.background = ctx.CreateLinearGradient(0, -d.h/2, 0, d.h/2)
	d.background.AddColorStopString(0, "black")
	d.background.AddColorStopString(0.5, "rgb(100, 0, 50)")
	d.background.AddColorStopString(1, "black")
	defer d.background.Release()

	ctx.Translate(d.w/2, d.h/2)
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
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (d *demo) update() {
	for i, p := range d.points {
		gone := false
		for j := range p {
			d.points[i][j].z -= 0.5
			if d.points[i][j].z < -300 {
				gone = true
			}
		}
		if gone {
			arr := d.points[len(d.points)-1]
			for k := range arr {
				dist := math.Abs(float64(k) - float64(len(arr))/2)
				arr[k].z = 0
				arr[k].y = rand.Float64()*-(dist*dist) + 30
			}
			copy(d.points[1:], d.points)
			d.points[0] = arr
		}
	}
}

func (d *demo) draw(ctx *canvas.Context) {
	ctx.SetFillStyleGradient(d.background)
	ctx.FillRect(-d.w/2, -d.h/2, d.w, d.h)
	ctx.BeginPath()
	ctx.Arc(0, 0, 200, 0, math.Pi*2, false)
	ctx.ClosePath()
	ctx.SetShadowColorString("orange")
	ctx.SetShadowBlur(100)
	ctx.SetFillStyleGradient(d.gradient)
	ctx.Fill()
	ctx.SetShadowBlur(0)
	for i := range len(d.points) - 1 {
		for j := range len(d.points[i]) - 1 {
			size := 300 / (300 + d.points[i][j].z)
			nextSize := 300 / (300 + d.points[i+1][j].z)
			ctx.BeginPath()
			ctx.MoveTo((d.points[i][j].x-d.offset)*size, d.points[i][j].y*size)
			ctx.LineTo((d.points[i][j+1].x-d.offset)*size, d.points[i][j+1].y*size)
			ctx.LineTo((d.points[i+1][j+1].x-d.offset)*nextSize, d.points[i+1][j+1].y*nextSize)
			ctx.LineTo((d.points[i+1][j].x-d.offset)*nextSize, d.points[i+1][j].y*nextSize)
			ctx.ClosePath()
			ctx.SetFillStyle(color.RGBA{A: uint8(limit(-d.points[i][j].z/100, 1.0) * 255)})
			c := 300 + d.points[i][j].z
			ctx.SetStrokeStyle(color.RGBA{
				R: uint8(250 - c),
				G: 0,
				B: uint8(50 + c),
				A: uint8((1 - c/300) * 255),
			})
			ctx.Fill()
			ctx.Stroke()
		}
	}
}

func limit(f, max float64) float64 {
	if f > max {
		return max
	}
	return f
}

func httpLink(addr string) string {
	if addr[0] == ':' {
		addr = "localhost" + addr
	}
	return "http://" + addr
}
