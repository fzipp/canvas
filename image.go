// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image"
	"image/draw"
)

type Image struct {
	id     uint32
	width  int
	height int
}

func (img *Image) Width() int {
	return img.width
}

func (img *Image) Height() int {
	return img.height
}

func ensureRGBA(img image.Image) *image.RGBA {
	switch im := img.(type) {
	case *image.RGBA:
		return im
	default:
		rgba := image.NewRGBA(im.Bounds())
		draw.Draw(rgba, im.Bounds(), im, image.Pt(0, 0), draw.Src)
		return rgba
	}
}
