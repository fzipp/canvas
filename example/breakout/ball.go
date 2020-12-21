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
	switch c {
	case collisionLeft, collisionRight:
		b.v.X = -b.v.X
	case collisionTop, collisionBottom:
		b.v.Y = -b.v.Y
	}
	return c
}

func (b *ball) checkCollision(rect image.Rectangle) collision {
	is := b.bounds().Intersect(rect)
	switch {
	case is.Max.X == rect.Max.X:
		return collisionRight
	case is.Min.X == rect.Min.X:
		return collisionLeft
	case is.Max.Y == rect.Max.Y:
		return collisionBottom
	case is.Min.Y == rect.Min.Y:
		return collisionTop
	}
	return collisionNone
}

type collision int

const (
	collisionNone = iota
	collisionLeft
	collisionRight
	collisionTop
	collisionBottom
)
