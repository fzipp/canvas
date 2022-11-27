// Copyright 2021 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image/color"
	"reflect"
	"testing"
)

func TestBufferWrite(t *testing.T) {
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

func TestBufferRead(t *testing.T) {
	tests := []struct {
		name      string
		bytes     []byte
		read      func(buf *buffer) any
		wantValue any
		wantBytes []byte
	}{
		{
			"readByte",
			[]byte{0x01, 0x02},
			func(buf *buffer) any {
				return buf.readByte()
			},
			byte(0x01),
			[]byte{0x02},
		},
		{
			"readUint32",
			[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			func(buf *buffer) any {
				return buf.readUint32()
			},
			uint32(0x01020304),
			[]byte{0x05, 0x06, 0x07, 0x08},
		},
		{
			"readUint64",
			[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09},
			func(buf *buffer) any {
				return buf.readUint64()
			},
			uint64(0x0102030405060708),
			[]byte{0x09},
		},
		{
			"readFloat64",
			[]byte{0x40, 0x09, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a, 0xab},
			func(buf *buffer) any {
				return buf.readFloat64()
			},
			3.2,
			[]byte{0xab},
		},
		{
			"readString",
			[]byte{0x00, 0x00, 0x00, 0x04, 0x54, 0x65, 0x73, 0x74, 0x42},
			func(buf *buffer) any {
				return buf.readString()
			},
			"Test",
			[]byte{0x42},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &buffer{bytes: tt.bytes}
			got := tt.read(buf)
			if buf.error != nil {
				t.Errorf("did not expect error, but got error: %s", buf.error)
			}
			if !reflect.DeepEqual(got, tt.wantValue) {
				t.Errorf("got: %#v, want: %#v", got, tt.wantValue)
			}
			if !reflect.DeepEqual(buf.bytes, tt.wantBytes) {
				t.Errorf("\n(buf.bytes) got : %#v\n(buf.bytes) want: %#v", buf.bytes, tt.wantBytes)
			}
		})
	}
}

func TestBufferReadErrors(t *testing.T) {
	tests := []struct {
		name      string
		bytes     []byte
		read      func(buf *buffer) any
		wantValue any
	}{
		{
			"readByte",
			[]byte{},
			func(buf *buffer) any {
				return buf.readByte()
			},
			byte(0),
		},
		{
			"readUint32",
			[]byte{0x01, 0x02, 0x03},
			func(buf *buffer) any {
				return buf.readUint32()
			},
			uint32(0),
		},
		{
			"readUint64",
			[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
			func(buf *buffer) any {
				return buf.readUint64()
			},
			uint64(0),
		},
		{
			"readFloat64",
			[]byte{0x40, 0x09, 0x99, 0x99, 0x99, 0x99, 0x99},
			func(buf *buffer) any {
				return buf.readFloat64()
			},
			float64(0),
		},
		{
			"readString: length data too short",
			[]byte{0x00, 0x00, 0x00},
			func(buf *buffer) any {
				return buf.readString()
			},
			"",
		},
		{
			"readString: string data too short",
			[]byte{0x00, 0x00, 0x00, 0x04, 0x54, 0x65, 0x73},
			func(buf *buffer) any {
				return buf.readString()
			},
			"",
		},
	}
	wantErrorMessage := "data too short"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &buffer{bytes: tt.bytes}
			got := tt.read(buf)
			if buf.error == nil {
				t.Errorf("expected error, but got none")
				return
			}
			if wantType, ok := buf.error.(errDataTooShort); !ok {
				t.Errorf("expected %T error, but got: %#v", wantType, buf.error)
			}
			if buf.error.Error() != wantErrorMessage {
				t.Errorf("expected %q error message, but got: %q", wantErrorMessage, buf.error)
			}
			if !reflect.DeepEqual(got, tt.wantValue) {
				t.Errorf("got: %#v, want: %#v", got, tt.wantValue)
			}
			if len(buf.bytes) != 0 {
				t.Errorf("excpected buf.bytes to be empty after short read, but got: %#v", buf.bytes)
			}
		})
	}
}
