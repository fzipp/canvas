// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image"
	"image/draw"
)

type ImageData struct {
	id     uint32
	ctx    *Context
	width  int
	height int
}

func (m *ImageData) Width() int {
	return m.width
}

func (m *ImageData) Height() int {
	return m.height
}

func (m *ImageData) Release() {
	m.ctx.buf.addByte(bReleaseImageData)
	m.ctx.buf.addUint32(m.id)
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
