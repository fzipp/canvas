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
		float64(b.radius), 0, 2*math.Pi, false)
	ctx.ClosePath()
	ctx.Fill()
}

func (b *ball) bounds() image.Rectangle {
	return image.Rect(
		int(b.pos.x-b.radius), int(b.pos.y-b.radius),
		int(b.pos.x+b.radius), int(b.pos.y+b.radius))
}

func (b *ball) bounceOnCollision(rect image.Rectangle) collision {
	c := b.checkCollision(rect)
	switch {
	case c&collisionTop > 0:
		b.pos.y = float64(rect.Min.Y) - b.radius - b.v.y
		b.v.y = -b.v.y
	case c&collisionBottom > 0:
		b.pos.y = float64(rect.Max.Y) + b.radius + b.v.y
		b.v.y = -b.v.y
	case c&collisionLeft > 0:
		b.pos.x = float64(rect.Min.X) - b.radius - b.v.x
		b.v.x = -b.v.x
	case c&collisionRight > 0:
		b.pos.x = float64(rect.Max.X) + b.radius + b.v.x
		b.v.x = -b.v.x
	}
	return c
}

func (b *ball) checkCollision(rect image.Rectangle) (c collision) {
	is := b.bounds().Intersect(rect)
	if is == (image.Rectangle{}) {
		return c
	}
	if is.Min.Y == rect.Min.Y {
		c |= collisionTop
	}
	if is.Max.Y == rect.Max.Y {
		c |= collisionBottom
	}
	if is.Min.X == rect.Min.X {
		c |= collisionLeft
	}
	if is.Max.X == rect.Max.X {
		c |= collisionRight
	}
	return c
}

type vec2 struct {
	x, y float64
}

func (v vec2) add(w vec2) vec2 {
	return vec2{x: v.x + w.x, y: v.y + w.y}
}

func (v vec2) sub(w vec2) vec2 {
	return vec2{x: v.x - w.x, y: v.y - w.y}
}

type collision int

const (
	collisionLeft collision = 1 << iota
	collisionRight
	collisionTop
	collisionBottom
	collisionNone collision = 0
)
