// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image"
	"image/color"
	"sync"
)

type Context struct {
	config config
	draws  chan<- []byte
	events <-chan Event
	buf    buffer

	imageDataIDs idGenerator
	gradientIDs  idGenerator
	patternIDs   idGenerator
}

func newContext(draws chan<- []byte, events <-chan Event, config config) *Context {
	return &Context{
		config: config,
		draws:  draws,
		events: events,
	}
}

func (ctx *Context) Events() <-chan Event {
	return ctx.events
}

func (ctx *Context) CanvasWidth() int {
	return ctx.config.width
}

func (ctx *Context) CanvasHeight() int {
	return ctx.config.height
}

func (ctx *Context) SetFillStyle(c color.Color) {
	ctx.buf.addByte(bFillStyle)
	ctx.buf.addColor(c)
}

func (ctx *Context) SetFillStyleString(color string) {
	ctx.buf.addByte(bFillStyleString)
	ctx.buf.addString(color)
}

func (ctx *Context) SetFillStyleGradient(g *Gradient) {
	ctx.buf.addByte(bFillStyleGradient)
	ctx.buf.addUint32(g.id)
}

func (ctx *Context) SetFillStylePattern(p *Pattern) {
	ctx.buf.addByte(bFillStylePattern)
	ctx.buf.addUint32(p.id)
}

func (ctx *Context) SetFont(font string) {
	ctx.buf.addByte(bFont)
	ctx.buf.addString(font)
}

func (ctx *Context) SetGlobalAlpha(alpha float64) {
	ctx.buf.addByte(bGlobalAlpha)
	ctx.buf.addFloat64(alpha)
}

func (ctx *Context) SetGlobalCompositeOperation(mode CompositeOperation) {
	ctx.buf.addByte(bGlobalCompositeOperation)
	ctx.buf.addByte(byte(mode))
}

func (ctx *Context) SetImageSmoothingEnabled(enabled bool) {
	ctx.buf.addByte(bImageSmoothingEnabled)
	ctx.buf.addBool(enabled)
}

func (ctx *Context) SetLineCap(cap LineCap) {
	ctx.buf.addByte(bLineCap)
	ctx.buf.addByte(byte(cap))
}

func (ctx *Context) SetLineDashOffset(offset float64) {
	ctx.buf.addByte(bLineDashOffset)
	ctx.buf.addFloat64(offset)
}

func (ctx *Context) SetLineJoin(join LineJoin) {
	ctx.buf.addByte(bLineJoin)
	ctx.buf.addByte(byte(join))
}

func (ctx *Context) SetLineWidth(width float64) {
	ctx.buf.addByte(bLineWidth)
	ctx.buf.addFloat64(width)
}

func (ctx *Context) SetMiterLimit(offset float64) {
	ctx.buf.addByte(bMiterLimit)
	ctx.buf.addFloat64(offset)
}

func (ctx *Context) SetShadowBlur(level float64) {
	ctx.buf.addByte(bShadowBlur)
	ctx.buf.addFloat64(level)
}

func (ctx *Context) SetShadowColor(c color.Color) {
	ctx.buf.addByte(bShadowColor)
	ctx.buf.addColor(c)
}

func (ctx *Context) SetShadowColorString(color string) {
	ctx.buf.addByte(bShadowColorString)
	ctx.buf.addString(color)
}

func (ctx *Context) SetShadowOffsetX(offset float64) {
	ctx.buf.addByte(bShadowOffsetX)
	ctx.buf.addFloat64(offset)
}

func (ctx *Context) SetShadowOffsetY(offset float64) {
	ctx.buf.addByte(bShadowOffsetY)
	ctx.buf.addFloat64(offset)
}

func (ctx *Context) SetStrokeStyle(c color.Color) {
	ctx.buf.addByte(bStrokeStyle)
	ctx.buf.addColor(c)
}

func (ctx *Context) SetStrokeStyleString(color string) {
	ctx.buf.addByte(bStrokeStyleString)
	ctx.buf.addString(color)
}

func (ctx *Context) SetStrokeStyleGradient(g *Gradient) {
	ctx.buf.addByte(bStrokeStyleGradient)
	ctx.buf.addUint32(g.id)
}

func (ctx *Context) SetStrokeStylePattern(p *Pattern) {
	ctx.buf.addByte(bStrokeStylePattern)
	ctx.buf.addUint32(p.id)
}

func (ctx *Context) SetTextAlign(align TextAlign) {
	ctx.buf.addByte(bTextAlign)
	ctx.buf.addByte(byte(align))
}

func (ctx *Context) SetTextBaseline(baseline TextBaseline) {
	ctx.buf.addByte(bTextBaseline)
	ctx.buf.addByte(byte(baseline))
}

func (ctx *Context) Arc(x, y, radius, startAngle, endAngle float64, anticlockwise bool) {
	ctx.buf.addByte(bArc)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(radius)
	ctx.buf.addFloat64(startAngle)
	ctx.buf.addFloat64(endAngle)
	ctx.buf.addBool(anticlockwise)
}

func (ctx *Context) ArcTo(x1, y1, x2, y2, radius float64) {
	ctx.buf.addByte(bArcTo)
	ctx.buf.addFloat64(x1)
	ctx.buf.addFloat64(y1)
	ctx.buf.addFloat64(x2)
	ctx.buf.addFloat64(y2)
	ctx.buf.addFloat64(radius)
}

func (ctx *Context) BeginPath() {
	ctx.buf.addByte(bBeginPath)
}

func (ctx *Context) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	ctx.buf.addByte(bBezierCurveTo)
	ctx.buf.addFloat64(cp1x)
	ctx.buf.addFloat64(cp1y)
	ctx.buf.addFloat64(cp2x)
	ctx.buf.addFloat64(cp2y)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

func (ctx *Context) ClearRect(x, y, width, height float64) {
	ctx.buf.addByte(bClearRect)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(width)
	ctx.buf.addFloat64(height)
}

func (ctx *Context) Clip() {
	ctx.buf.addByte(bClip)
}

func (ctx *Context) ClosePath() {
	ctx.buf.addByte(bClosePath)
}

func (ctx *Context) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, anticlockwise bool) {
	ctx.buf.addByte(bEllipse)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(radiusX)
	ctx.buf.addFloat64(radiusY)
	ctx.buf.addFloat64(rotation)
	ctx.buf.addFloat64(startAngle)
	ctx.buf.addFloat64(endAngle)
	ctx.buf.addBool(anticlockwise)
}

func (ctx *Context) Fill() {
	ctx.buf.addByte(bFill)
}

func (ctx *Context) FillRect(x, y, width, height float64) {
	ctx.buf.addByte(bFillRect)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(width)
	ctx.buf.addFloat64(height)
}

func (ctx *Context) FillText(text string, x, y float64) {
	ctx.buf.addByte(bFillText)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addString(text)
}

func (ctx *Context) FillTextMaxWidth(text string, x, y, maxWidth float64) {
	ctx.buf.addByte(bFillTextMaxWidth)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(maxWidth)
	ctx.buf.addString(text)
}

func (ctx *Context) LineTo(x, y float64) {
	ctx.buf.addByte(bLineTo)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

func (ctx *Context) MoveTo(x, y float64) {
	ctx.buf.addByte(bMoveTo)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

func (ctx *Context) QuadraticCurveTo(cpx, cpy, x, y float64) {
	ctx.buf.addByte(bQuadraticCurveTo)
	ctx.buf.addFloat64(cpx)
	ctx.buf.addFloat64(cpy)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

func (ctx *Context) Rect(x, y, width, height float64) {
	ctx.buf.addByte(bRect)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(width)
	ctx.buf.addFloat64(height)
}

func (ctx *Context) Restore() {
	ctx.buf.addByte(bRestore)
}

func (ctx *Context) Rotate(angle float64) {
	ctx.buf.addByte(bRotate)
	ctx.buf.addFloat64(angle)
}

func (ctx *Context) Save() {
	ctx.buf.addByte(bSave)
}

func (ctx *Context) Scale(x, y float64) {
	ctx.buf.addByte(bScale)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

func (ctx *Context) Stroke() {
	ctx.buf.addByte(bStroke)
}

func (ctx *Context) StrokeText(text string, x, y float64) {
	ctx.buf.addByte(bStrokeText)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addString(text)
}

func (ctx *Context) StrokeTextMaxWidth(text string, x, y, maxWidth float64) {
	ctx.buf.addByte(bStrokeTextMaxWidth)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(maxWidth)
	ctx.buf.addString(text)
}

func (ctx *Context) StrokeRect(x, y, width, height float64) {
	ctx.buf.addByte(bStrokeRect)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(width)
	ctx.buf.addFloat64(height)
}

func (ctx *Context) Translate(x, y float64) {
	ctx.buf.addByte(bTranslate)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

func (ctx *Context) Transform(a, b, c, d, e, f float64) {
	ctx.buf.addByte(bTransform)
	ctx.buf.addFloat64(a)
	ctx.buf.addFloat64(b)
	ctx.buf.addFloat64(c)
	ctx.buf.addFloat64(d)
	ctx.buf.addFloat64(e)
	ctx.buf.addFloat64(f)
}

func (ctx *Context) SetTransform(a, b, c, d, e, f float64) {
	ctx.buf.addByte(bSetTransform)
	ctx.buf.addFloat64(a)
	ctx.buf.addFloat64(b)
	ctx.buf.addFloat64(c)
	ctx.buf.addFloat64(d)
	ctx.buf.addFloat64(e)
	ctx.buf.addFloat64(f)
}

func (ctx *Context) SetLineDash(segments []float64) {
	ctx.buf.addByte(bSetLineDash)
	ctx.buf.addUint32(uint32(len(segments)))
	for _, seg := range segments {
		ctx.buf.addFloat64(seg)
	}
}

func (ctx *Context) CreateImageData(m image.Image) *ImageData {
	rgba := ensureRGBA(m)
	bounds := m.Bounds()
	id := ctx.imageDataIDs.GenerateID()
	ctx.buf.addByte(bCreateImageData)
	ctx.buf.addUint32(id)
	ctx.buf.addUint32(uint32(bounds.Dx()))
	ctx.buf.addUint32(uint32(bounds.Dy()))
	ctx.buf.addBytes(rgba.Pix)
	return &ImageData{id: id, ctx: ctx, width: bounds.Dx(), height: bounds.Dy()}
}

func (ctx *Context) PutImageData(src *ImageData, dx, dy float64) {
	ctx.buf.addByte(bPutImageData)
	ctx.buf.addUint32(src.id)
	ctx.buf.addFloat64(dx)
	ctx.buf.addFloat64(dy)
}

func (ctx *Context) PutImageDataDirty(src *ImageData, dx, dy, dirtyX, dirtyY, dirtyWidth, dirtyHeight float64) {
	ctx.buf.addByte(bPutImageDataDirty)
	ctx.buf.addUint32(src.id)
	ctx.buf.addFloat64(dx)
	ctx.buf.addFloat64(dy)
	ctx.buf.addFloat64(dirtyX)
	ctx.buf.addFloat64(dirtyY)
	ctx.buf.addFloat64(dirtyWidth)
	ctx.buf.addFloat64(dirtyHeight)
}

func (ctx *Context) DrawImage(src *ImageData, dx, dy float64) {
	ctx.buf.addByte(bDrawImage)
	ctx.buf.addUint32(src.id)
	ctx.buf.addFloat64(dx)
	ctx.buf.addFloat64(dy)
}

func (ctx *Context) DrawImageScaled(src *ImageData, dx, dy, dWidth, dHeight float64) {
	ctx.buf.addByte(bDrawImageScaled)
	ctx.buf.addUint32(src.id)
	ctx.buf.addFloat64(dx)
	ctx.buf.addFloat64(dy)
	ctx.buf.addFloat64(dWidth)
	ctx.buf.addFloat64(dHeight)
}

func (ctx *Context) DrawImageSubRectangle(src *ImageData, sx, sy, sWidth, sHeight, dx, dy, dWidth, dHeight float64) {
	ctx.buf.addByte(bDrawImageSubRectangle)
	ctx.buf.addUint32(src.id)
	ctx.buf.addFloat64(sx)
	ctx.buf.addFloat64(sy)
	ctx.buf.addFloat64(sWidth)
	ctx.buf.addFloat64(sHeight)
	ctx.buf.addFloat64(dx)
	ctx.buf.addFloat64(dy)
	ctx.buf.addFloat64(dWidth)
	ctx.buf.addFloat64(dHeight)
}

func (ctx *Context) CreateLinearGradient(x0, y0, x1, y1 float64) *Gradient {
	id := ctx.gradientIDs.GenerateID()
	ctx.buf.addByte(bCreateLinearGradient)
	ctx.buf.addUint32(id)
	ctx.buf.addFloat64(x0)
	ctx.buf.addFloat64(y0)
	ctx.buf.addFloat64(x1)
	ctx.buf.addFloat64(y1)
	return &Gradient{id: id, ctx: ctx}
}

func (ctx *Context) CreateRadialGradient(x0, y0, r0, x1, y1, r1 float64) *Gradient {
	id := ctx.gradientIDs.GenerateID()
	ctx.buf.addByte(bCreateRadialGradient)
	ctx.buf.addUint32(id)
	ctx.buf.addFloat64(x0)
	ctx.buf.addFloat64(y0)
	ctx.buf.addFloat64(r0)
	ctx.buf.addFloat64(x1)
	ctx.buf.addFloat64(y1)
	ctx.buf.addFloat64(r1)
	return &Gradient{id: id, ctx: ctx}
}

func (ctx *Context) CreatePattern(src *ImageData, repetition PatternRepetition) *Pattern {
	id := ctx.patternIDs.GenerateID()
	ctx.buf.addByte(bCreatePattern)
	ctx.buf.addUint32(id)
	ctx.buf.addUint32(src.id)
	ctx.buf.addByte(byte(repetition))
	return &Pattern{id: id, ctx: ctx}
}

func (ctx *Context) GetImageData(sx, sy, sw, sh float64) *ImageData {
	id := ctx.imageDataIDs.GenerateID()
	ctx.buf.addByte(bGetImageData)
	ctx.buf.addUint32(id)
	ctx.buf.addFloat64(sx)
	ctx.buf.addFloat64(sy)
	ctx.buf.addFloat64(sw)
	ctx.buf.addFloat64(sh)
	return &ImageData{id: id, ctx: ctx, width: int(sw), height: int(sh)}
}

func (ctx *Context) Flush() {
	ctx.draws <- ctx.buf.bytes
	ctx.buf.reset()
}

type idGenerator struct {
	nextMu sync.Mutex
	next   uint32
}

func (g *idGenerator) GenerateID() uint32 {
	g.nextMu.Lock()
	defer g.nextMu.Unlock()
	id := g.next
	g.next++
	return id
}
