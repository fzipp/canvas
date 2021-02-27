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
	pos    image.Point
	radius int
	v      image.Point
	color  color.Color
}

func (b *ball) update() {
	b.pos = b.pos.Add(b.v)
}

func (b *ball) draw(ctx *canvas.Context) {
	ctx.SetFillStyle(b.color)
	ctx.BeginPath()
	ctx.Arc(
		float64(b.pos.X), float64(b.pos.Y),
		float64(b.radius), 0, 2*math.Pi, false)
	ctx.ClosePath()
	ctx.Fill()
}

func (b *ball) bounds() image.Rectangle {
	return image.Rect(
		b.pos.X-b.radius, b.pos.Y-b.radius,
		b.pos.X+b.radius, b.pos.Y+b.radius)
}

func (b *ball) bounceOnCollision(rect image.Rectangle) collision {
	c := b.checkCollision(rect)
	switch {
	case c&collisionTop > 0:
		b.pos.Y = rect.Min.Y - b.radius - b.v.Y
		b.v.Y = -b.v.Y
	case c&collisionBottom > 0:
		b.pos.Y = rect.Max.Y + b.radius + b.v.Y
		b.v.Y = -b.v.Y
	case c&collisionLeft > 0:
		b.pos.X = rect.Min.X - b.radius - b.v.X
		b.v.X = -b.v.X
	case c&collisionRight > 0:
		b.pos.X = rect.Max.X + b.radius + b.v.X
		b.v.X = -b.v.X
	}
	return c
}

func (b *ball) checkCollision(rect image.Rectangle) collision {
	is := b.bounds().Intersect(rect)
	c := collisionNone
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

type collision int

const (
	collisionLeft collision = 1 << iota
	collisionRight
	collisionTop
	collisionBottom
	collisionNone collision = 0
)
