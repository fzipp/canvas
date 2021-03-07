// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"image/color"
	"math"

	"github.com/fzipp/canvas"
)

type ball struct {
	pos    vec2
	radius float64
	v      vec2
	color  color.Color
}

func (b *ball) update() {
	b.pos = b.pos.add(b.v)
}

func (b *ball) draw(ctx *canvas.Context) {
	ctx.SetFillStyle(b.color)
	ctx.BeginPath()
	ctx.Arc(
		b.pos.x, b.pos.y,
		b.radius, 0, 2*math.Pi, false)
	ctx.ClosePath()
	ctx.Fill()
}

func (b *ball) bounds() image.Rectangle {
	return image.Rect(
		int(b.pos.x-b.radius), int(b.pos.y-b.radius),
		int(b.pos.x+b.radius), int(b.pos.y+b.radius))
}

func (b *ball) bounceOnCollision(rect image.Rectangle) bool {
	n := b.checkCollision(rect)
	if n == (vec2{}) {
		return false
	}
	b.pos = b.pos.sub(b.v)
	b.v = b.v.reflect(n)
	b.pos = b.pos.add(b.v)
	return true
}

func (b *ball) checkCollision(rect image.Rectangle) (normal vec2) {
	is := b.bounds().Intersect(rect)
	if is == (image.Rectangle{}) {
		return normal
	}
	if is.Min.Y == rect.Min.Y {
		normal = normal.add(vec2{0, -1})
	}
	if is.Max.Y == rect.Max.Y {
		normal = normal.add(vec2{0, 1})
	}
	if is.Min.X == rect.Min.X {
		normal = normal.add(vec2{-1, 0})
	}
	if is.Max.X == rect.Max.X {
		normal = normal.add(vec2{1, 0})
	}
	return normal.norm()
}
