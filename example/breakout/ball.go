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

func (b *ball) bounceOnCollision(rect image.Rectangle) Collision {
	c := b.checkCollision(rect)
	switch c {
	case CollisionLeft, CollisionRight:
		b.v.X = -b.v.X
	case CollisionTop, CollisionBottom:
		b.v.Y = -b.v.Y
	}
	return c
}

func (b *ball) checkCollision(rect image.Rectangle) Collision {
	is := b.bounds().Intersect(rect)
	switch {
	case is.Max.X == rect.Max.X:
		return CollisionRight
	case is.Min.X == rect.Min.X:
		return CollisionLeft
	case is.Max.Y == rect.Max.Y:
		return CollisionBottom
	case is.Min.Y == rect.Min.Y:
		return CollisionTop
	}
	return CollisionNone
}

type Collision int

const (
	CollisionNone = iota
	CollisionLeft
	CollisionRight
	CollisionTop
	CollisionBottom
)
