// Copyright 2021 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image/color"
	"reflect"
	"testing"
)

func TestBuffer(t *testing.T) {
	tests := []struct {
		name string
		use  func(buf *buffer)
		want []byte
	}{
		{
			"addByte",
			func(buf *buffer) {
				buf.addByte(0xfa)
				buf.addByte(0x12)
				buf.addByte(0x01)
			},
			[]byte{0xfa, 0x12, 0x01},
		},
		{
			"addFloat64",
			func(buf *buffer) {
				buf.addFloat64(3.1415)
				buf.addFloat64(362.5)
			},
			[]byte{
				0x40, 0x09, 0x21, 0xca, 0xc0, 0x83, 0x12, 0x6f,
				0x40, 0x76, 0xa8, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"addUint32",
			func(buf *buffer) {
				buf.addUint32(12)
				buf.addUint32(4096)
				buf.addUint32(60000)
				buf.addUint32(1200000)
			},
			[]byte{
				0x00, 0x00, 0x00, 0x0c,
				0x00, 0x00, 0x10, 0x00,
				0x00, 0x00, 0xea, 0x60,
				0x00, 0x12, 0x4f, 0x80,
			},
		},
		{
			"addBool",
			func(buf *buffer) {
				buf.addBool(true)
				buf.addBool(false)
				buf.addBool(true)
			},
			[]byte{0x01, 0x00, 0x01},
		},
		{
			"addBytes",
			func(buf *buffer) {
				buf.addBytes([]byte{0x01, 0x02, 0xbc, 0xbd})
				buf.addBytes([]byte{0xfe, 0xff})
			},
			[]byte{0x01, 0x02, 0xbc, 0xbd, 0xfe, 0xff},
		},
		{
			"addString",
			func(buf *buffer) {
				buf.addString("hello")
				buf.addString("äöü")
			},
			[]byte{
				0x00, 0x00, 0x00, 0x05, // len(s)
				0x68, 0x65, 0x6c, 0x6c, 0x6f,
				0x00, 0x00, 0x00, 0x06, // len(s)
				0xc3, 0xa4, 0xc3, 0xb6, 0xc3, 0xbc,
			},
		},
		{
			"addColor",
			func(buf *buffer) {
				buf.addColor(color.Black)
				buf.addColor(color.White)
				buf.addColor(color.RGBA{R: 127, G: 32, B: 64, A: 255})
			},
			[]byte{
				0x00, 0x00, 0x00, 0xff,
				0xff, 0xff, 0xff, 0xff,
				0x7f, 0x20, 0x40, 0xff,
			},
		},
		{
			"reset",
			func(buf *buffer) {
				buf.addBytes([]byte{0x01, 0x02, 0x03})
				buf.reset()
				buf.addBytes([]byte{0x04, 0x05, 0x06, 0x07})
			},
			[]byte{0x04, 0x05, 0x06, 0x07},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &buffer{}
			tt.use(buf)
			got := buf.bytes
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot : %#v\nwant: %#v", got, tt.want)
			}
		})
	}
}
