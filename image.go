// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The API doc comments are based on the MDN Web Docs for the [Canvas API]
// by Mozilla Contributors and are licensed under [CC-BY-SA 2.5].
//
// [Canvas API]: https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D
// [CC-BY-SA 2.5]: https://creativecommons.org/licenses/by-sa/2.5/

package canvas

import (
	"image"
	"image/draw"
)

// ImageData represents the underlying pixel data of an image. It is
// created using the Context.CreateImageData and Context.GetImageData methods.
// It can also be used to set a part of the canvas by using
// Context.PutImageData, Context.PutImageDataDirty, Context.DrawImage and
// Context.DrawImageScaled.
//
// The image data should be released with the Release method when it is no
// longer needed.
type ImageData struct {
	id       uint32
	ctx      *Context
	width    int
	height   int
	released bool
}

// Width returns the actual width, in pixels, of the image.
func (m *ImageData) Width() int {
	return m.width
}

// Height returns the actual height, in pixels, of the image.
func (m *ImageData) Height() int {
	return m.height
}

// Release releases the image data on the client side.
func (m *ImageData) Release() {
	if m.released {
		return
	}
	m.ctx.buf.addByte(bReleaseImageData)
	m.ctx.buf.addUint32(m.id)
	m.released = true
}

func (m *ImageData) checkUseAfterRelease() {
	if m.released {
		panic("ImageData: use after release")
	}
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
