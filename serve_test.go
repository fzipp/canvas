// Copyright 2022 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image/color"
	"testing"
)

func TestRGBAString(t *testing.T) {
	tests := []struct {
		name  string
		color color.Color
		want  string
	}{
		{"zero", color.RGBA{}, "rgba(0, 0, 0, 0)"},
		{"black", color.Black, "rgba(0, 0, 0, 1)"},
		{"white", color.White, "rgba(255, 255, 255, 1)"},
		{"red", color.RGBA{R: 0xFF, A: 0xFF}, "rgba(255, 0, 0, 1)"},
		{"green", color.RGBA{G: 0xFF, A: 0xFF}, "rgba(0, 255, 0, 1)"},
		{"blue", color.RGBA{B: 0xFF, A: 0xFF}, "rgba(0, 0, 255, 1)"},
		{"grey", color.RGBA{R: 0x7F, G: 0x7F, B: 0x7F, A: 0xFF}, "rgba(127, 127, 127, 1)"},
		{"grey semi-transparent", color.RGBA{R: 0x7F, G: 0x7F, B: 0x7F, A: 0x7F}, "rgba(127, 127, 127, 0.5)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rgbaString(tt.color)
			if got != tt.want {
				t.Errorf("rgbaString(%v) = %v, want: %v", tt.color, got, tt.want)
			}
		})
	}
}
