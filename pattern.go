// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

// Pattern represents a pattern, based on an image, created by the
// Context.CreatePattern method. It can be used with the
// Context.SetFillStylePattern and Context.SetStrokeStylePattern methods.
//
// The pattern should be released with the Release method when it is no longer
// needed.
type Pattern struct {
	id  uint32
	ctx *Context
}

// Release releases the pattern on the client side.
func (p *Pattern) Release() {
	p.ctx.buf.addByte(bReleasePattern)
	p.ctx.buf.addUint32(p.id)
}
