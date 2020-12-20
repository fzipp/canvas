// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

type Pattern struct {
	id  uint32
	ctx *Context
}

func (p *Pattern) Release() {
	msg := [1 + 4]byte{bReleasePattern}
	byteOrder.PutUint32(msg[1:], p.id)
	p.ctx.write(msg[:])
}
