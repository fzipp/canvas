// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The API doc comments are based on the MDN Web Docs for the [Canvas API]
// by Mozilla Contributors and are licensed under [CC-BY-SA 2.5].
//
// [Canvas API]: https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D
// [CC-BY-SA 2.5]: https://creativecommons.org/licenses/by-sa/2.5/

package canvas

// Pattern represents a pattern, based on an image, created by the
// Context.CreatePattern method. It can be used with the
// Context.SetFillStylePattern and Context.SetStrokeStylePattern methods.
//
// The pattern should be released with the Release method when it is no longer
// needed.
type Pattern struct {
	id       uint32
	ctx      *Context
	released bool
}

// Release releases the pattern on the client side.
func (p *Pattern) Release() {
	if p.released {
		return
	}
	p.ctx.buf.addByte(bReleasePattern)
	p.ctx.buf.addUint32(p.id)
	p.released = true
}

func (p *Pattern) checkUseAfterRelease() {
	if p.released {
		panic("Pattern: use after release")
	}
}
