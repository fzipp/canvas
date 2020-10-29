// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"image/color"

	"github.com/fzipp/canvas"
)

type paddle struct {
	pos   image.Point
	size  image.Point
	color color.Color
}

func (p *paddle) draw(ctx *canvas.Context) {
	ctx.SetFillStyle(p.color)
	x := p.pos.X - (p.size.X / 2)
	y := p.pos.Y - (p.size.Y / 2)
	ctx.FillRect(float64(x), float64(y), float64(p.size.X), float64(p.size.Y))
}

func (p *paddle) bounds() image.Rectangle {
	return image.Rect(
		p.pos.X-p.size.X/2, p.pos.Y-p.size.Y/2,
		p.pos.X+p.size.X/2, p.pos.Y+p.size.Y/2)
}
