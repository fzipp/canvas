// Copyright 2021 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/fzipp/canvas"
)

func main() {
	port := ":8080"
	fmt.Println("Listening on http://localhost" + port)
	err := canvas.ListenAndServe(port, run,
		canvas.Title("Hilbert"),
		canvas.Size(1000, 500),
		canvas.FullPage(),
	)
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
