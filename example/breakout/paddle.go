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
	pos   vec2
	size  vec2
	color color.Color
}

func (p *paddle) draw(ctx *canvas.Context) {
	ctx.SetFillStyle(p.color)
	x := p.pos.x - (p.size.x / 2)
	y := p.pos.y - (p.size.y / 2)
	ctx.FillRect(x, y, p.size.x, p.size.y)
}

func (p *paddle) bounds() image.Rectangle {
	return image.Rect(
		int(p.pos.x-p.size.x/2), int(p.pos.y-p.size.y/2),
		int(p.pos.x+p.size.x/2), int(p.pos.y+p.size.y/2))
}
