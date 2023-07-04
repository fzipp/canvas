// Copyright 2021 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hilbert draws a graphic pattern consisting of multiple superimposed Hilbert
// curves using the recursive algorithm described in the book "Algorithms and
// Data Structures" by N. Wirth.
//
// Usage:
//
//	hilbert [-http address]
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

	"github.com/fzipp/canvas"
)

func main() {
	http := flag.String("http", ":8080", "HTTP service address (e.g., '127.0.0.1:8080' or just ':8080')")
	flag.Parse()

	fmt.Println("Listening on " + httpLink(*http))
	err := canvas.ListenAndServe(*http, run, &canvas.Options{
		Title:             "Hilbert",
		Width:             500,
		Height:            500,
		ScaleToPageHeight: true,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *canvas.Context) {
	ctx.SetLineWidth(2)
	ctx.SetStrokeStyleString("#777")
	h := &hilbert{ctx: ctx}
	h.draw()
	ctx.Flush()
}

type hilbert struct {
	ctx     *canvas.Context
	x, y, d int
}

func (h *hilbert) draw() {
	fw := h.ctx.CanvasWidth()
	fh := h.ctx.CanvasHeight()
	w := fw
	if w > fh {
		w = fh
	}
	h.d = 8
	k := 0
	for h.d*2 < w {
		h.d *= 2
		k++
	}
	x0 := fw / 2
	y0 := fh / 2
	for n := 1; n <= k; n++ {
		h.d /= 2
		x0 += h.d / 2
		y0 -= h.d / 2
		h.x = x0
		h.y = y0
		h.ctx.MoveTo(float64(h.x), float64(h.y))
		h.A(n)
	}
	h.ctx.Stroke()
}

func (h *hilbert) A(i int) {
	if i > 0 {
		h.D(i - 1)
		h.W()
		h.A(i - 1)
		h.S()
		h.A(i - 1)
		h.E()
		h.B(i - 1)
	}
}

func (h *hilbert) B(i int) {
	if i > 0 {
		h.C(i - 1)
		h.N()
		h.B(i - 1)
		h.E()
		h.B(i - 1)
		h.S()
		h.A(i - 1)
	}
}

func (h *hilbert) C(i int) {
	if i > 0 {
		h.B(i - 1)
		h.E()
		h.C(i - 1)
		h.N()
		h.C(i - 1)
		h.W()
		h.D(i - 1)
	}
}

func (h *hilbert) D(i int) {
	if i > 0 {
		h.A(i - 1)
		h.S()
		h.D(i - 1)
		h.W()
		h.D(i - 1)
		h.N()
		h.C(i - 1)
	}
}

func (h *hilbert) E() {
	h.x += h.d
	h.ctx.LineTo(float64(h.x), float64(h.y))
}

func (h *hilbert) N() {
	h.y -= h.d
	h.ctx.LineTo(float64(h.x), float64(h.y))
}

func (h *hilbert) W() {
	h.x -= h.d
	h.ctx.LineTo(float64(h.x), float64(h.y))
}

func (h *hilbert) S() {
	h.y += h.d
	h.ctx.LineTo(float64(h.x), float64(h.y))
}

func httpLink(addr string) string {
	if addr[0] == ':' {
		addr = "localhost" + addr
	}
	return "http://" + addr
}
