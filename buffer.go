// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"encoding/binary"
	"errors"
	"image/color"
	"math"
)

type buffer struct {
	bytes []byte
	error error
}

var byteOrder = binary.BigEndian

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

func (buf *buffer) readByte() byte {
	if len(buf.bytes) < 1 {
		buf.dataTooShort()
		return 0
	}
	b := buf.bytes[0]
	buf.bytes = buf.bytes[1:]
	return b
}

func (buf *buffer) readUint32() uint32 {
	if len(buf.bytes) < 4 {
		buf.dataTooShort()
		return 0
	}
	i := byteOrder.Uint32(buf.bytes)
	buf.bytes = buf.bytes[4:]
	return i
}

func (buf *buffer) readUint64() uint64 {
	if len(buf.bytes) < 8 {
		buf.dataTooShort()
		return 0
	}
	i := byteOrder.Uint64(buf.bytes)
	buf.bytes = buf.bytes[8:]
	return i
}

func (buf *buffer) readFloat64() float64 {
	return math.Float64frombits(buf.readUint64())
}

func (buf *buffer) readString() string {
	length := int(buf.readUint32())
	if len(buf.bytes) < length {
		buf.dataTooShort()
		return ""
	}
	s := string(buf.bytes[:length])
	buf.bytes = buf.bytes[length:]
	return s
}

func (buf *buffer) reset() {
	buf.bytes = make([]byte, 0, cap(buf.bytes))
}

func (buf *buffer) dataTooShort() {
	buf.reset()
	buf.error = errors.New("data too short")
}
