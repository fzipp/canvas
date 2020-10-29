// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"image/color"

	"github.com/fzipp/canvas"
)

type brick struct {
	rect   image.Rectangle
	color  color.Color
	points int
}

func (b *brick) draw(ctx *canvas.Context) {
	ctx.SetFillStyle(b.color)
	ctx.FillRect(
		float64(b.rect.Min.X), float64(b.rect.Min.Y),
		float64(b.rect.Dx()), float64(b.rect.Dy()))
}

func (b *brick) bounds() image.Rectangle {
	return b.rect
}
