// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import "image/color"

// Gradient represents a gradient. It is returned by the methods
// Context.CreateLinearGradient and Context.CreateRadialGradient.
// It can be used with the Context.SetFillStyleGradient and
// Context.SetStrokeStyleGradient methods.
//
// The gradient should be released with the Release method when it is no longer
// needed.
type Gradient struct {
	id       uint32
	ctx      *Context
	released bool
}

// AddColorStop adds a new stop, defined by an offset and a color, to the
// gradient.
func (g *Gradient) AddColorStop(offset float64, c color.Color) {
	g.checkUseAfterRelease()
	g.ctx.buf.addByte(bGradientAddColorStop)
	g.ctx.buf.addUint32(g.id)
	g.ctx.buf.addFloat64(offset)
	g.ctx.buf.addColor(c)
}

// AddColorStopString adds a new stop, defined by an offset and a color, to
// the gradient.
//
// The color is parsed as a CSS color value like "#a100cb", "#ccc",
// "darkgreen", "rgba(0.5, 0.2, 0.7, 1.0)", etc.
func (g *Gradient) AddColorStopString(offset float64, color string) {
	g.checkUseAfterRelease()
	g.ctx.buf.addByte(bGradientAddColorStopString)
	g.ctx.buf.addUint32(g.id)
	g.ctx.buf.addFloat64(offset)
	g.ctx.buf.addString(color)
}

// Release releases the gradient on the client side.
func (g *Gradient) Release() {
	if g.released {
		return
	}
	g.ctx.buf.addByte(bReleaseGradient)
	g.ctx.buf.addUint32(g.id)
	g.released = true
}

func (g *Gradient) checkUseAfterRelease() {
	if g.released {
		panic("Gradient: use after release")
	}
}
