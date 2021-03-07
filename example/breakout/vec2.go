// Copyright 2021 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "math"

type vec2 struct {
	x, y float64
}

func (v vec2) add(w vec2) vec2 {
	return vec2{x: v.x + w.x, y: v.y + w.y}
}

func (v vec2) sub(w vec2) vec2 {
	return vec2{x: v.x - w.x, y: v.y - w.y}
}

func (v vec2) dot(w vec2) float64 {
	return v.x*w.x + v.y*w.y
}

func (v vec2) mul(s float64) vec2 {
	return vec2{v.x * s, v.y * s}
}

func (v vec2) div(s float64) vec2 {
	return vec2{v.x / s, v.y / s}
}

func (v vec2) norm() vec2 {
	return v.div(v.len())
}

func (v vec2) len() float64 {
	return math.Sqrt(v.sqLen())
}

func (v vec2) sqLen() float64 {
	return v.dot(v)
}

func (v vec2) reflect(n vec2) vec2 {
	return v.sub(n.mul(2 * v.dot(n)))
}
