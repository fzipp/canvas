// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"math"
)

var byteOrder = binary.BigEndian

type Context struct {
	config config
	draws  chan<- []byte
	events <-chan Event
	quit   <-chan struct{}
	buf    bytes.Buffer

	nextImageID    uint32
	nextGradientID uint32
}

func newContext(draws chan<- []byte, events <-chan Event, quit <-chan struct{}, config config) *Context {
	return &Context{
		config: config,
		draws:  draws,
		events: events,
		quit:   quit,
	}
}

func (ctx *Context) Events() <-chan Event {
	return ctx.events
}

func (ctx *Context) Quit() <-chan struct{} {
	return ctx.quit
}

func (ctx *Context) CanvasWidth() int {
	return ctx.config.width
}

func (ctx *Context) CanvasHeight() int {
	return ctx.config.height
}

func (ctx *Context) SetFillStyle(c color.Color) {
	clr := color.RGBAModel.Convert(c).(color.RGBA)
	ctx.write([]byte{bFillStyle, clr.R, clr.G, clr.B, clr.A})
}

func (ctx *Context) SetFillStyleString(color string) {
	msg := make([]byte, 1+4+len(color))
	msg[0] = bFillStyleString
	byteOrder.PutUint32(msg[1:], uint32(len(color)))
	copy(msg[5:], color)
	ctx.write(msg)
}

func (ctx *Context) SetFillStyleGradient(g *Gradient) {
	msg := [1 + 4]byte{bFillStyleGradient}
	byteOrder.PutUint32(msg[1:], g.id)
	ctx.write(msg[:])
}

func (ctx *Context) SetFont(font string) {
	msg := make([]byte, 1+4+len(font))
	msg[0] = bFont
	byteOrder.PutUint32(msg[1:], uint32(len(font)))
	copy(msg[5:], font)
	ctx.write(msg)
}

func (ctx *Context) SetGlobalAlpha(alpha float64) {
	msg := [1 + 8]byte{bGlobalAlpha}
	byteOrder.PutUint64(msg[1:], math.Float64bits(alpha))
	ctx.write(msg[:])
}

func (ctx *Context) SetGlobalCompositeOperation(mode CompositeOperation) {
	ctx.write([]byte{bGlobalCompositeOperation, byte(mode)})
}

func (ctx *Context) SetImageSmoothingEnabled(enabled bool) {
	msg := []byte{bImageSmoothingEnabled, 0}
	if enabled {
		msg[1] = 1
	}
	ctx.write(msg)
}

func (ctx *Context) SetLineCap(cap LineCap) {
	ctx.write([]byte{bLineCap, byte(cap)})
}

func (ctx *Context) SetLineDashOffset(offset float64) {
	msg := [1 + 8]byte{bLineDashOffset}
	byteOrder.PutUint64(msg[1:], math.Float64bits(offset))
	ctx.write(msg[:])
}

func (ctx *Context) SetLineJoin(join LineJoin) {
	ctx.write([]byte{bLineJoin, byte(join)})
}

func (ctx *Context) SetLineWidth(width float64) {
	msg := [1 + 8]byte{bLineWidth}
	byteOrder.PutUint64(msg[1:], math.Float64bits(width))
	ctx.write(msg[:])
}

func (ctx *Context) SetMiterLimit(offset float64) {
	msg := [1 + 8]byte{bMiterLimit}
	byteOrder.PutUint64(msg[1:], math.Float64bits(offset))
	ctx.write(msg[:])
}

func (ctx *Context) SetShadowBlur(level float64) {
	msg := [1 + 8]byte{bShadowBlur}
	byteOrder.PutUint64(msg[1:], math.Float64bits(level))
	ctx.write(msg[:])
}

func (ctx *Context) SetShadowColor(c color.Color) {
	clr := color.RGBAModel.Convert(c).(color.RGBA)
	ctx.write([]byte{bShadowColor, clr.R, clr.G, clr.B, clr.A})
}

func (ctx *Context) SetShadowColorString(color string) {
	msg := make([]byte, 1+4+len(color))
	msg[0] = bShadowColorString
	byteOrder.PutUint32(msg[1:], uint32(len(color)))
	copy(msg[5:], color)
	ctx.write(msg)
}

func (ctx *Context) SetShadowOffsetX(offset float64) {
	msg := [1 + 8]byte{bShadowOffsetX}
	byteOrder.PutUint64(msg[1:], math.Float64bits(offset))
	ctx.write(msg[:])
}

func (ctx *Context) SetShadowOffsetY(offset float64) {
	msg := [1 + 8]byte{bShadowOffsetY}
	byteOrder.PutUint64(msg[1:], math.Float64bits(offset))
	ctx.write(msg[:])
}

func (ctx *Context) SetStrokeStyle(c color.Color) {
	clr := color.RGBAModel.Convert(c).(color.RGBA)
	ctx.write([]byte{bStrokeStyle, clr.R, clr.G, clr.B, clr.A})
}

func (ctx *Context) SetStrokeStyleString(color string) {
	msg := make([]byte, 1+4+len(color))
	msg[0] = bStrokeStyleString
	byteOrder.PutUint32(msg[1:], uint32(len(color)))
	copy(msg[5:], color)
	ctx.write(msg)
}

func (ctx *Context) SetStrokeStyleGradient(g *Gradient) {
	msg := [1 + 4]byte{bStrokeStyleGradient}
	byteOrder.PutUint32(msg[1:], g.id)
	ctx.write(msg[:])
}

func (ctx *Context) SetTextAlign(align TextAlign) {
	ctx.write([]byte{bTextAlign, byte(align)})
}

func (ctx *Context) SetTextBaseline(baseline TextBaseline) {
	ctx.write([]byte{bTextBaseline, byte(baseline)})
}

func (ctx *Context) Arc(x, y, radius, startAngle, endAngle float64, anticlockwise bool) {
	msg := [1 + 5*8 + 1]byte{bArc}
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	byteOrder.PutUint64(msg[17:], math.Float64bits(radius))
	byteOrder.PutUint64(msg[25:], math.Float64bits(startAngle))
	byteOrder.PutUint64(msg[33:], math.Float64bits(endAngle))
	if anticlockwise {
		msg[41] = 1
	}
	ctx.write(msg[:])
}

func (ctx *Context) ArcTo(x1, y1, x2, y2, radius float64) {
	msg := [1 + 5*8]byte{bArcTo}
	byteOrder.PutUint64(msg[1:], math.Float64bits(x1))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y1))
	byteOrder.PutUint64(msg[17:], math.Float64bits(x2))
	byteOrder.PutUint64(msg[25:], math.Float64bits(y2))
	byteOrder.PutUint64(msg[33:], math.Float64bits(radius))
	ctx.write(msg[:])
}

func (ctx *Context) BeginPath() {
	ctx.write([]byte{bBeginPath})
}

func (ctx *Context) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	msg := [1 + 6*8]byte{bBezierCurveTo}
	byteOrder.PutUint64(msg[1:], math.Float64bits(cp1x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(cp1y))
	byteOrder.PutUint64(msg[17:], math.Float64bits(cp2x))
	byteOrder.PutUint64(msg[25:], math.Float64bits(cp2y))
	byteOrder.PutUint64(msg[33:], math.Float64bits(x))
	byteOrder.PutUint64(msg[41:], math.Float64bits(y))
	ctx.write(msg[:])
}

func (ctx *Context) ClearRect(x, y, width, height float64) {
	msg := [1 + 4*8]byte{bClearRect}
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	byteOrder.PutUint64(msg[17:], math.Float64bits(width))
	byteOrder.PutUint64(msg[25:], math.Float64bits(height))
	ctx.write(msg[:])
}

func (ctx *Context) Clip() {
	ctx.write([]byte{bClip})
}

func (ctx *Context) ClosePath() {
	ctx.write([]byte{bClosePath})
}

func (ctx *Context) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, anticlockwise bool) {
	msg := [1 + 7*8 + 1]byte{bEllipse}
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	byteOrder.PutUint64(msg[17:], math.Float64bits(radiusX))
	byteOrder.PutUint64(msg[25:], math.Float64bits(radiusY))
	byteOrder.PutUint64(msg[33:], math.Float64bits(rotation))
	byteOrder.PutUint64(msg[41:], math.Float64bits(startAngle))
	byteOrder.PutUint64(msg[49:], math.Float64bits(endAngle))
	if anticlockwise {
		msg[57] = 1
	}
	ctx.write(msg[:])
}

func (ctx *Context) Fill() {
	ctx.write([]byte{bFill})
}

func (ctx *Context) FillRect(x, y, width, height float64) {
	msg := [1 + 4*8]byte{bFillRect}
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	byteOrder.PutUint64(msg[17:], math.Float64bits(width))
	byteOrder.PutUint64(msg[25:], math.Float64bits(height))
	ctx.write(msg[:])
}

func (ctx *Context) FillText(text string, x, y float64) {
	msg := make([]byte, 1+2*8+4+len(text))
	msg[0] = bFillText
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	byteOrder.PutUint32(msg[17:], uint32(len(text)))
	copy(msg[21:], text)
	ctx.write(msg)
}

func (ctx *Context) FillTextMaxWidth(text string, x, y, maxWidth float64) {
	msg := make([]byte, 1+3*8+4+len(text))
	msg[0] = bFillTextMaxWidth
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	byteOrder.PutUint64(msg[17:], math.Float64bits(maxWidth))
	byteOrder.PutUint32(msg[25:], uint32(len(text)))
	copy(msg[29:], text)
	ctx.write(msg)
}

func (ctx *Context) LineTo(x, y float64) {
	msg := [1 + 2*8]byte{bLineTo}
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	ctx.write(msg[:])
}

func (ctx *Context) MoveTo(x, y float64) {
	msg := [1 + 2*8]byte{bMoveTo}
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	ctx.write(msg[:])
}

func (ctx *Context) QuadraticCurveTo(cpx, cpy, x, y float64) {
	msg := [1 + 4*8]byte{bQuadraticCurveTo}
	byteOrder.PutUint64(msg[1:], math.Float64bits(cpx))
	byteOrder.PutUint64(msg[9:], math.Float64bits(cpy))
	byteOrder.PutUint64(msg[17:], math.Float64bits(x))
	byteOrder.PutUint64(msg[25:], math.Float64bits(y))
	ctx.write(msg[:])
}

func (ctx *Context) Rect(x, y, width, height float64) {
	msg := [1 + 4*8]byte{bRect}
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	byteOrder.PutUint64(msg[17:], math.Float64bits(width))
	byteOrder.PutUint64(msg[25:], math.Float64bits(height))
	ctx.write(msg[:])
}

func (ctx *Context) Restore() {
	ctx.write([]byte{bRestore})
}

func (ctx *Context) Rotate(angle float64) {
	msg := [1 + 8]byte{bRotate}
	byteOrder.PutUint64(msg[1:], math.Float64bits(angle))
	ctx.write(msg[:])
}

func (ctx *Context) Save() {
	ctx.write([]byte{bSave})
}

func (ctx *Context) Scale(x, y float64) {
	msg := [1 + 2*8]byte{bScale}
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	ctx.write(msg[:])
}

func (ctx *Context) Stroke() {
	ctx.write([]byte{bStroke})
}

func (ctx *Context) StrokeText(text string, x, y float64) {
	msg := make([]byte, 1+2*8+4+len(text))
	msg[0] = bStrokeText
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	byteOrder.PutUint32(msg[17:], uint32(len(text)))
	copy(msg[21:], text)
	ctx.write(msg)
}

func (ctx *Context) StrokeTextMaxWidth(text string, x, y, maxWidth float64) {
	msg := make([]byte, 1+3*8+4+len(text))
	msg[0] = bStrokeTextMaxWidth
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	byteOrder.PutUint64(msg[17:], math.Float64bits(maxWidth))
	byteOrder.PutUint32(msg[25:], uint32(len(text)))
	copy(msg[29:], text)
	ctx.write(msg)
}

func (ctx *Context) StrokeRect(x, y, width, height float64) {
	msg := [1 + 4*8]byte{bStrokeRect}
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	byteOrder.PutUint64(msg[17:], math.Float64bits(width))
	byteOrder.PutUint64(msg[25:], math.Float64bits(height))
	ctx.write(msg[:])
}

func (ctx *Context) Translate(x, y float64) {
	msg := [1 + 2*8]byte{bTranslate}
	byteOrder.PutUint64(msg[1:], math.Float64bits(x))
	byteOrder.PutUint64(msg[9:], math.Float64bits(y))
	ctx.write(msg[:])
}

func (ctx *Context) Transform(a, b, c, d, e, f float64) {
	msg := [1 + 6*8]byte{bTransform}
	byteOrder.PutUint64(msg[1:], math.Float64bits(a))
	byteOrder.PutUint64(msg[9:], math.Float64bits(b))
	byteOrder.PutUint64(msg[17:], math.Float64bits(c))
	byteOrder.PutUint64(msg[25:], math.Float64bits(d))
	byteOrder.PutUint64(msg[33:], math.Float64bits(e))
	byteOrder.PutUint64(msg[41:], math.Float64bits(f))
	ctx.write(msg[:])
}

func (ctx *Context) SetTransform(a, b, c, d, e, f float64) {
	msg := [1 + 6*8]byte{bSetTransform}
	byteOrder.PutUint64(msg[1:], math.Float64bits(a))
	byteOrder.PutUint64(msg[9:], math.Float64bits(b))
	byteOrder.PutUint64(msg[17:], math.Float64bits(c))
	byteOrder.PutUint64(msg[25:], math.Float64bits(d))
	byteOrder.PutUint64(msg[33:], math.Float64bits(e))
	byteOrder.PutUint64(msg[41:], math.Float64bits(f))
	ctx.write(msg[:])
}

func (ctx *Context) SetLineDash(segments []float64) {
	msg := make([]byte, 1+4+len(segments)*8)
	msg[0] = bSetLineDash
	byteOrder.PutUint32(msg[1:], uint32(len(segments)))
	for i, seg := range segments {
		byteOrder.PutUint64(msg[5+i*8:], math.Float64bits(seg))
	}
	ctx.write(msg[:])
}

func (ctx *Context) CreateImageData(img image.Image) *Image {
	rgba := ensureRGBA(img)
	msg := make([]byte, 1+3*4+len(rgba.Pix))
	msg[0] = bCreateImageData
	bounds := img.Bounds()
	id := ctx.nextImageID
	ctx.nextImageID++
	byteOrder.PutUint32(msg[1:], id)
	byteOrder.PutUint32(msg[5:], uint32(bounds.Dx()))
	byteOrder.PutUint32(msg[9:], uint32(bounds.Dy()))
	copy(msg[13:], rgba.Pix)
	ctx.write(msg)
	return &Image{id: id, ctx: ctx, width: bounds.Dx(), height: bounds.Dy()}
}

func (ctx *Context) PutImageData(img *Image, dx, dy float64) {
	msg := [1 + 4 + 2*8]byte{bPutImageData}
	byteOrder.PutUint32(msg[1:], img.id)
	byteOrder.PutUint64(msg[5:], math.Float64bits(dx))
	byteOrder.PutUint64(msg[13:], math.Float64bits(dy))
	ctx.write(msg[:])
}

func (ctx *Context) PutImageDataDirty(img *Image, dx, dy, dirtyX, dirtyY, dirtyWidth, dirtyHeight float64) {
	msg := [1 + 4 + 6*8]byte{bPutImageDataDirty}
	byteOrder.PutUint32(msg[1:], img.id)
	byteOrder.PutUint64(msg[5:], math.Float64bits(dx))
	byteOrder.PutUint64(msg[13:], math.Float64bits(dy))
	byteOrder.PutUint64(msg[21:], math.Float64bits(dirtyX))
	byteOrder.PutUint64(msg[29:], math.Float64bits(dirtyY))
	byteOrder.PutUint64(msg[37:], math.Float64bits(dirtyWidth))
	byteOrder.PutUint64(msg[45:], math.Float64bits(dirtyHeight))
	ctx.write(msg[:])
}

func (ctx *Context) DrawImage(img *Image, dx, dy float64) {
	msg := [1 + 4 + 2*8]byte{bDrawImage}
	byteOrder.PutUint32(msg[1:], img.id)
	byteOrder.PutUint64(msg[5:], math.Float64bits(dx))
	byteOrder.PutUint64(msg[13:], math.Float64bits(dy))
	ctx.write(msg[:])
}

func (ctx *Context) DrawImageScaled(img *Image, dx, dy, dWidth, dHeight float64) {
	msg := [1 + 4 + 4*8]byte{bDrawImageScaled}
	byteOrder.PutUint32(msg[1:], img.id)
	byteOrder.PutUint64(msg[5:], math.Float64bits(dx))
	byteOrder.PutUint64(msg[13:], math.Float64bits(dy))
	byteOrder.PutUint64(msg[21:], math.Float64bits(dWidth))
	byteOrder.PutUint64(msg[29:], math.Float64bits(dHeight))
	ctx.write(msg[:])
}

func (ctx *Context) DrawImageSubRectangle(img *Image, sx, sy, sWidth, sHeight, dx, dy, dWidth, dHeight float64) {
	msg := [1 + 4 + 8*8]byte{bDrawImageSubRectangle}
	byteOrder.PutUint32(msg[1:], img.id)
	byteOrder.PutUint64(msg[5:], math.Float64bits(sx))
	byteOrder.PutUint64(msg[13:], math.Float64bits(sy))
	byteOrder.PutUint64(msg[21:], math.Float64bits(sWidth))
	byteOrder.PutUint64(msg[29:], math.Float64bits(sHeight))
	byteOrder.PutUint64(msg[37:], math.Float64bits(dx))
	byteOrder.PutUint64(msg[45:], math.Float64bits(dy))
	byteOrder.PutUint64(msg[53:], math.Float64bits(dWidth))
	byteOrder.PutUint64(msg[61:], math.Float64bits(dHeight))
	ctx.write(msg[:])
}

func (ctx *Context) CreateLinearGradient(x0, y0, x1, y1 float64) *Gradient {
	id := ctx.nextGradientID
	ctx.nextGradientID++
	msg := [1 + 4 + 4*8]byte{bCreateLinearGradient}
	byteOrder.PutUint32(msg[1:], id)
	byteOrder.PutUint64(msg[5:], math.Float64bits(x0))
	byteOrder.PutUint64(msg[13:], math.Float64bits(y0))
	byteOrder.PutUint64(msg[21:], math.Float64bits(x1))
	byteOrder.PutUint64(msg[29:], math.Float64bits(y1))
	ctx.write(msg[:])
	return &Gradient{id: id, ctx: ctx}
}

func (ctx *Context) CreateRadialGradient(x0, y0, r0, x1, y1, r1 float64) *Gradient {
	id := ctx.nextGradientID
	ctx.nextGradientID++
	msg := [1 + 4 + 6*8]byte{bCreateRadialGradient}
	byteOrder.PutUint32(msg[1:], id)
	byteOrder.PutUint64(msg[5:], math.Float64bits(x0))
	byteOrder.PutUint64(msg[13:], math.Float64bits(y0))
	byteOrder.PutUint64(msg[21:], math.Float64bits(r0))
	byteOrder.PutUint64(msg[29:], math.Float64bits(x1))
	byteOrder.PutUint64(msg[37:], math.Float64bits(y1))
	byteOrder.PutUint64(msg[45:], math.Float64bits(r1))
	ctx.write(msg[:])
	return &Gradient{id: id, ctx: ctx}
}

func (ctx *Context) Flush() {
	select {
	case <-ctx.Quit():
		return
	case ctx.draws <- ctx.buf.Bytes():
		ctx.buf = bytes.Buffer{}
	}
}

func (ctx *Context) write(p []byte) {
	ctx.buf.Write(p)
}
