// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import "image/color"

type Gradient struct {
	id  uint32
	ctx *Context
}

func (g *Gradient) AddColorStop(offset float64, c color.Color) {
	g.ctx.buf.addByte(bGradientAddColorStop)
	g.ctx.buf.addUint32(g.id)
	g.ctx.buf.addFloat64(offset)
	g.ctx.buf.addColor(c)
}

func (g *Gradient) AddColorStopString(offset float64, color string) {
	g.ctx.buf.addByte(bGradientAddColorStopString)
	g.ctx.buf.addUint32(g.id)
	g.ctx.buf.addFloat64(offset)
	g.ctx.buf.addString(color)
}

func (g *Gradient) Release() {
	g.ctx.buf.addByte(bReleaseGradient)
	g.ctx.buf.addUint32(g.id)
}
