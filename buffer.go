// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image/color"
	"math"
)

type buffer struct {
	bytes []byte
}

func (buf *buffer) addByte(b byte) {
	buf.bytes = append(buf.bytes, b)
}

func (buf *buffer) addFloat64(f float64) {
	buf.bytes = append(buf.bytes, 0, 0, 0, 0, 0, 0, 0, 0)
	byteOrder.PutUint64(buf.bytes[len(buf.bytes)-8:], math.Float64bits(f))
}

func (buf *buffer) addUint32(i uint32) {
	buf.bytes = append(buf.bytes, 0, 0, 0, 0)
	byteOrder.PutUint32(buf.bytes[len(buf.bytes)-4:], i)
}

func (buf *buffer) addBool(b bool) {
	if b {
		buf.addByte(1)
	} else {
		buf.addByte(0)
	}
}

func (buf *buffer) addBytes(p []byte) {
	buf.bytes = append(buf.bytes, p...)
}

func (buf *buffer) addString(s string) {
	buf.addUint32(uint32(len(s)))
	buf.bytes = append(buf.bytes, []byte(s)...)
}

func (buf *buffer) addColor(c color.Color) {
	clr := color.RGBAModel.Convert(c).(color.RGBA)
	buf.addByte(clr.R)
	buf.addByte(clr.G)
	buf.addByte(clr.B)
	buf.addByte(clr.A)
}

func (buf *buffer) reset() {
	buf.bytes = buf.bytes[:0]
}
