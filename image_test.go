// Copyright 2021 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func TestEnsureRGBA(t *testing.T) {
	tests := []struct {
		name string
		img  image.Image
		want *image.RGBA
	}{
		{
			"Gray",
			&image.Gray{
				Pix: []uint8{
					0x00, 0x4a,
					0xb0, 0x12,
					0xff, 0xa0,
				},
				Stride: 2,
				Rect:   image.Rect(0, 0, 2, 3),
			},
			&image.RGBA{
				Pix: []uint8{
					0x00, 0x00, 0x00, 0xff, 0x4a, 0x4a, 0x4a, 0xff,
					0xb0, 0xb0, 0xb0, 0xff, 0x12, 0x12, 0x12, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xa0, 0xa0, 0xa0, 0xff,
				},
				Stride: 8,
				Rect:   image.Rect(0, 0, 2, 3),
			},
		},
		{
			"Paletted",
			&image.Paletted{
				Pix: []uint8{
					0x01, 0x02,
					0x03, 0x01,
					0x00, 0x02,
				},
				Stride: 2,
				Rect:   image.Rect(0, 0, 2, 3),
				Palette: color.Palette{
					color.RGBA{R: 0xa1, G: 0xa2, B: 0xa3, A: 0xa4},
					color.RGBA{R: 0xb1, G: 0xb2, B: 0xb3, A: 0xb4},
					color.RGBA{R: 0xc1, G: 0xc2, B: 0xc3, A: 0xc4},
					color.RGBA{R: 0xd1, G: 0xd2, B: 0xd3, A: 0xd4},
				},
			},
			&image.RGBA{
				Pix: []uint8{
					0xb1, 0xb2, 0xb3, 0xb4, 0xc1, 0xc2, 0xc3, 0xc4,
					0xd1, 0xd2, 0xd3, 0xd4, 0xb1, 0xb2, 0xb3, 0xb4,
					0xa1, 0xa2, 0xa3, 0xa4, 0xc1, 0xc2, 0xc3, 0xc4,
				},
				Stride: 8,
				Rect:   image.Rect(0, 0, 2, 3),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ensureRGBA(tt.img)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\n got: %#v\nwant: %#v", got, tt.want)
			}
		})
	}
}
