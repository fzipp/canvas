// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image/color"
	"math"
)

type Gradient struct {
	id  uint32
	ctx *Context
}

func (g *Gradient) AddColorStop(offset float64, c color.Color) {
	clr := color.RGBAModel.Convert(c).(color.RGBA)
	msg := [1 + 4 + 8 + 4]byte{bGradientAddColorStop}
	byteOrder.PutUint32(msg[1:], g.id)
	byteOrder.PutUint64(msg[5:], math.Float64bits(offset))
	msg[13] = clr.R
	msg[14] = clr.G
	msg[15] = clr.B
	msg[16] = clr.A
	g.ctx.write(msg[:])
}

func (g *Gradient) AddColorStopString(offset float64, color string) {
	msg := make([]byte, 1+4+8+4+len(color))
	msg[0] = bGradientAddColorStopString
	byteOrder.PutUint32(msg[1:], g.id)
	byteOrder.PutUint64(msg[5:], math.Float64bits(offset))
	byteOrder.PutUint32(msg[13:], uint32(len(color)))
	copy(msg[17:], color)
	g.ctx.write(msg)
}

func (g *Gradient) Release() {
	msg := [1 + 4]byte{bReleaseGradient}
	byteOrder.PutUint32(msg[1:], g.id)
	g.ctx.write(msg[:])
}
