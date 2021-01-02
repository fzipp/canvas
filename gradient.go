// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image/color"
)

type Gradient struct {
	id  uint32
	ctx *Context
}

func (g *Gradient) AddColorStop(offset float64, c color.Color) {
	buf := g.ctx.buf
	buf.addByte(bGradientAddColorStop)
	buf.addUint32(g.id)
	buf.addFloat64(offset)
	buf.addColor(c)
}

func (g *Gradient) AddColorStopString(offset float64, color string) {
	buf := g.ctx.buf
	buf.addByte(bGradientAddColorStopString)
	buf.addUint32(g.id)
	buf.addFloat64(offset)
	buf.addString(color)
}

func (g *Gradient) Release() {
	buf := g.ctx.buf
	buf.addByte(bReleaseGradient)
	buf.addUint32(g.id)
}
