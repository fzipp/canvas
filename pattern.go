// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

type Pattern struct {
	id  uint32
	ctx *Context
}

func (p *Pattern) Release() {
	p.ctx.buf.addByte(bReleasePattern)
	p.ctx.buf.addUint32(p.id)
}
