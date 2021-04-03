// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image"
	"image/color"
	"math"
	"reflect"
	"testing"
)

func TestContextDrawing(t *testing.T) {
	tests := []struct {
		name string
		draw func(*Context)
		want []byte
	}{
		{
			"Arc",
			func(ctx *Context) { ctx.Arc(5, 10, 15, 0, 0.75, true) },
			[]byte{
				0x01,
				0x40, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x2e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x3f, 0xe8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x01,
			},
		},
		{
			"ArcTo",
			func(ctx *Context) { ctx.ArcTo(100, 50, 80, 60, 25) },
			[]byte{
				0x02,
				0x40, 0x59, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x49, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x4e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x39, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"BeginPath",
			func(ctx *Context) { ctx.BeginPath() },
			[]byte{0x03},
		},
		{
			"BezierCurveTo",
			func(ctx *Context) { ctx.BezierCurveTo(230, 30, 150, 80, 250, 100) },
			[]byte{
				0x04,
				0x40, 0x6c, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x62, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x6f, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x59, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"ClearRect",
			func(ctx *Context) { ctx.ClearRect(10, 20, 120, 80) },
			[]byte{
				0x05,
				0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x5e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"Clip",
			func(ctx *Context) { ctx.Clip() },
			[]byte{0x06},
		},
		{
			"ClosePath",
			func(ctx *Context) { ctx.ClosePath() },
			[]byte{0x07},
		},
		{
			"CreateImageData",
			func(ctx *Context) {
				img := image.NewRGBA(image.Rect(0, 0, 4, 3))
				for i := range img.Pix {
					img.Pix[i] = byte(i)
				}
				ctx.CreateImageData(img)
			},
			[]byte{
				0x08,
				0x00, 0x00, 0x00, 0x00, // ID
				0x00, 0x00, 0x00, 0x04, // Width
				0x00, 0x00, 0x00, 0x03, // Height
				// Pixels
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
			},
		},
		{
			"CreateImageData: ID generation",
			func(ctx *Context) {
				ctx.CreateImageData(image.NewRGBA(image.Rect(0, 0, 0, 0)))
				ctx.CreateImageData(image.NewRGBA(image.Rect(0, 0, 0, 0)))
				ctx.CreateImageData(image.NewRGBA(image.Rect(0, 0, 0, 0)))
			},
			[]byte{
				0x08,                   // CreateImageData
				0x00, 0x00, 0x00, 0x00, // ID = 0
				0x00, 0x00, 0x00, 0x00, // Width
				0x00, 0x00, 0x00, 0x00, // Height
				0x08,                   // CreateImageData
				0x00, 0x00, 0x00, 0x01, // ID = 1
				0x00, 0x00, 0x00, 0x00, // Width
				0x00, 0x00, 0x00, 0x00, // Height
				0x08,                   // CreateImageData
				0x00, 0x00, 0x00, 0x02, // ID = 2
				0x00, 0x00, 0x00, 0x00, // Width
				0x00, 0x00, 0x00, 0x00, // Height
			},
		},
		{
			"Ellipse",
			func(ctx *Context) { ctx.Ellipse(100, 80, 50, 75, math.Pi/4, 0, 2*math.Pi, true) },
			[]byte{
				0x0e,
				0x40, 0x59, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x49, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x52, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x3f, 0xe9, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x19, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18,
				0x01,
			},
		},
		{
			"Fill",
			func(ctx *Context) { ctx.Fill() },
			[]byte{0x0f},
		},
		{
			"FillRect",
			func(ctx *Context) { ctx.FillRect(250, 120, 70, 65) },
			[]byte{
				0x10,
				0x40, 0x6f, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x5e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x51, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x50, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"SetFillStyle",
			func(ctx *Context) { ctx.SetFillStyle(color.RGBA{R: 10, G: 20, B: 30, A: 255}) },
			[]byte{
				0x11,
				0x0a, 0x14, 0x1e, 0xff,
			},
		},
		{
			"SetFillStyleString",
			func(ctx *Context) { ctx.SetFillStyleString("#0A141E") },
			[]byte{
				0x3b,
				0x00, 0x00, 0x00, 0x07, // len(color)
				0x23, 0x30, 0x41, 0x31, 0x34, 0x31, 0x45, // color
			},
		},
		{
			"SetFillStyleGradient",
			func(ctx *Context) { ctx.SetFillStyleGradient(&Gradient{id: 5}) },
			[]byte{
				0x16,
				0x00, 0x00, 0x00, 0x05,
			},
		},
		{
			"SetFillStylePattern",
			func(ctx *Context) { ctx.SetFillStylePattern(&Pattern{id: 16}) },
			[]byte{
				0x42,
				0x00, 0x00, 0x00, 0x10,
			},
		},
		{
			"FillText",
			func(ctx *Context) { ctx.FillText("test äöü", 22, 38) },
			[]byte{
				0x12,
				0x40, 0x36, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // x
				0x40, 0x43, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y
				0x00, 0x00, 0x00, 0x0b, // len(text)
				0x74, 0x65, 0x73, 0x74, 0x20, 0xc3, 0xa4, 0xc3, 0xb6, 0xc3, 0xbc, // text
			},
		},
		{
			"FillTextMaxWidth",
			func(ctx *Context) { ctx.FillTextMaxWidth("Hello, 世界", 45, 52, 100) },
			[]byte{
				0x39,
				0x40, 0x46, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, // x
				0x40, 0x4a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y
				0x40, 0x59, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // maxWidth
				0x00, 0x00, 0x00, 0x0d, // len(text)
				0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, // text "Hello, "
				0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c, // text "世界"
			},
		},
		{
			"SetFont",
			func(ctx *Context) { ctx.SetFont("Helvetica") },
			[]byte{
				0x13,
				0x00, 0x00, 0x00, 0x09, // len(font)
				0x48, 0x65, 0x6c, 0x76, 0x65, 0x74, 0x69, 0x63, 0x61, // font
			},
		},
		{
			"SetGlobalAlpha",
			func(ctx *Context) { ctx.SetGlobalAlpha(0.5) },
			[]byte{
				0x17,
				0x3f, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"SetGlobalCompositeOperation",
			func(ctx *Context) {
				ctx.SetGlobalCompositeOperation(OpSourceOver)
				ctx.SetGlobalCompositeOperation(OpSourceIn)
				ctx.SetGlobalCompositeOperation(OpSourceOut)
				ctx.SetGlobalCompositeOperation(OpSourceAtop)
				ctx.SetGlobalCompositeOperation(OpDestinationOver)
				ctx.SetGlobalCompositeOperation(OpDestinationIn)
				ctx.SetGlobalCompositeOperation(OpDestinationOut)
				ctx.SetGlobalCompositeOperation(OpDestinationAtop)
				ctx.SetGlobalCompositeOperation(OpLighter)
				ctx.SetGlobalCompositeOperation(OpCopy)
				ctx.SetGlobalCompositeOperation(OpXOR)
				ctx.SetGlobalCompositeOperation(OpMultiply)
				ctx.SetGlobalCompositeOperation(OpScreen)
				ctx.SetGlobalCompositeOperation(OpOverlay)
				ctx.SetGlobalCompositeOperation(OpDarken)
				ctx.SetGlobalCompositeOperation(OpLighten)
				ctx.SetGlobalCompositeOperation(OpColorDodge)
				ctx.SetGlobalCompositeOperation(OpColorBurn)
				ctx.SetGlobalCompositeOperation(OpHardLight)
				ctx.SetGlobalCompositeOperation(OpSoftLight)
				ctx.SetGlobalCompositeOperation(OpDifference)
				ctx.SetGlobalCompositeOperation(OpExclusion)
				ctx.SetGlobalCompositeOperation(OpHue)
				ctx.SetGlobalCompositeOperation(OpSaturation)
				ctx.SetGlobalCompositeOperation(OpColor)
				ctx.SetGlobalCompositeOperation(OpLuminosity)
			},
			[]byte{
				0x18, 0x00,
				0x18, 0x01,
				0x18, 0x02,
				0x18, 0x03,
				0x18, 0x04,
				0x18, 0x05,
				0x18, 0x06,
				0x18, 0x07,
				0x18, 0x08,
				0x18, 0x09,
				0x18, 0x0a,
				0x18, 0x0b,
				0x18, 0x0c,
				0x18, 0x0d,
				0x18, 0x0e,
				0x18, 0x0f,
				0x18, 0x10,
				0x18, 0x11,
				0x18, 0x12,
				0x18, 0x13,
				0x18, 0x14,
				0x18, 0x15,
				0x18, 0x16,
				0x18, 0x17,
				0x18, 0x18,
				0x18, 0x19,
			},
		},
		{
			"SetImageSmoothingEnabled",
			func(ctx *Context) {
				ctx.SetImageSmoothingEnabled(false)
				ctx.SetImageSmoothingEnabled(true)
			},
			[]byte{
				0x19, 0x00,
				0x19, 0x01,
			},
		},
		{
			"SetLineCap",
			func(ctx *Context) {
				ctx.SetLineCap(CapButt)
				ctx.SetLineCap(CapRound)
				ctx.SetLineCap(CapSquare)
			},
			[]byte{
				0x1c, 0x00,
				0x1c, 0x01,
				0x1c, 0x02,
			},
		},
		{
			"SetLineDashOffset",
			func(ctx *Context) { ctx.SetLineDashOffset(0.3) },
			[]byte{
				0x1d,
				0x3f, 0xd3, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
			},
		},
		{
			"SetLineJoin",
			func(ctx *Context) {
				ctx.SetLineJoin(JoinMiter)
				ctx.SetLineJoin(JoinRound)
				ctx.SetLineJoin(JoinBevel)
			},
			[]byte{
				0x1e, 0x00,
				0x1e, 0x01,
				0x1e, 0x02,
			},
		},
		{
			"LineTo",
			func(ctx *Context) { ctx.LineTo(10, 20) },
			[]byte{
				0x1f,
				0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"SetLineWidth",
			func(ctx *Context) { ctx.SetLineWidth(2.5) },
			[]byte{
				0x20,
				0x40, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"SetMiterLimit",
			func(ctx *Context) { ctx.SetMiterLimit(5) },
			[]byte{
				0x22,
				0x40, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"SetShadowBlur",
			func(ctx *Context) { ctx.SetShadowBlur(15) },
			[]byte{
				0x2d,
				0x40, 0x2e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"SetShadowColor",
			func(ctx *Context) { ctx.SetShadowColor(color.RGBA{R: 10, G: 20, B: 30, A: 255}) },
			[]byte{
				0x2e,
				0x0a, 0x14, 0x1e, 0xff,
			},
		},
		{
			"SetShadowColorString",
			func(ctx *Context) { ctx.SetShadowColorString("#a0b3ff") },
			[]byte{
				0x3d,
				0x00, 0x00, 0x00, 0x07, // len(color)
				0x23, 0x61, 0x30, 0x62, 0x33, 0x66, 0x66, // color
			},
		},
		{
			"SetShadowOffsetX",
			func(ctx *Context) { ctx.SetShadowOffsetX(32) },
			[]byte{
				0x2f,
				0x40, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"SetShadowOffsetY",
			func(ctx *Context) { ctx.SetShadowOffsetY(25) },
			[]byte{
				0x30,
				0x40, 0x39, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"SetStrokeStyle",
			func(ctx *Context) { ctx.SetStrokeStyle(color.RGBA{R: 120, G: 80, B: 110, A: 255}) },
			[]byte{
				0x33,
				0x78, 0x50, 0x6e, 0xff,
			},
		},
		{
			"SetStrokeStyleString",
			func(ctx *Context) { ctx.SetStrokeStyleString("LightGreen") },
			[]byte{
				0x3c,
				0x00, 0x00, 0x00, 0x0a, // len(color)
				0x4c, 0x69, 0x67, 0x68, 0x74, 0x47, 0x72, 0x65, 0x65, 0x6e, // color
			},
		},
		{
			"SetStrokeStyleGradient",
			func(ctx *Context) { ctx.SetStrokeStyleGradient(&Gradient{id: 333333}) },
			[]byte{
				0x1a,
				0x00, 0x05, 0x16, 0x15,
			},
		},
		{
			"SetStrokeStylePattern",
			func(ctx *Context) { ctx.SetStrokeStylePattern(&Pattern{id: 54321}) },
			[]byte{
				0x43,
				0x00, 0x00, 0xd4, 0x31,
			},
		},
		{
			"SetTextAlign",
			func(ctx *Context) {
				ctx.SetTextAlign(AlignStart)
				ctx.SetTextAlign(AlignEnd)
				ctx.SetTextAlign(AlignLeft)
				ctx.SetTextAlign(AlignRight)
				ctx.SetTextAlign(AlignCenter)
			},
			[]byte{
				0x35, 0x00,
				0x35, 0x01,
				0x35, 0x02,
				0x35, 0x03,
				0x35, 0x04,
			},
		},
		{
			"SetTextBaseline",
			func(ctx *Context) {
				ctx.SetTextBaseline(BaselineAlphabetic)
				ctx.SetTextBaseline(BaselineIdeographic)
				ctx.SetTextBaseline(BaselineTop)
				ctx.SetTextBaseline(BaselineBottom)
				ctx.SetTextBaseline(BaselineHanging)
				ctx.SetTextBaseline(BaselineMiddle)
			},
			[]byte{
				0x36, 0x00,
				0x36, 0x01,
				0x36, 0x02,
				0x36, 0x03,
				0x36, 0x04,
				0x36, 0x05,
			},
		},
		{
			"MoveTo",
			func(ctx *Context) { ctx.MoveTo(300, 200) },
			[]byte{
				0x23,
				0x40, 0x72, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x69, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"QuadraticCurveTo",
			func(ctx *Context) { ctx.QuadraticCurveTo(130, 215, 155, 330) },
			[]byte{
				0x25,
				0x40, 0x60, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x6a, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x63, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x74, 0xa0, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"Rect",
			func(ctx *Context) { ctx.Rect(52.5, 82.5, 70.2, 120.8) },
			[]byte{
				0x26,
				0x40, 0x4a, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x54, 0xa0, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x51, 0x8c, 0xcc, 0xcc, 0xcc, 0xcc, 0xcd,
				0x40, 0x5e, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
			},
		},
		{
			"Restore",
			func(ctx *Context) { ctx.Restore() },
			[]byte{0x27},
		},
		{
			"Rotate",
			func(ctx *Context) { ctx.Rotate(3.1415) },
			[]byte{
				0x28,
				0x40, 0x09, 0x21, 0xca, 0xc0, 0x83, 0x12, 0x6f,
			},
		},
		{
			"Save",
			func(ctx *Context) { ctx.Save() },
			[]byte{0x29},
		},
		{
			"Scale",
			func(ctx *Context) { ctx.Scale(2.5, 3) },
			[]byte{
				0x2a,
				0x40, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"Stroke",
			func(ctx *Context) { ctx.Stroke() },
			[]byte{0x31},
		},
		{
			"StrokeText",
			func(ctx *Context) { ctx.StrokeText("Test", 32, 42) },
			[]byte{
				0x34,
				0x40, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // x
				0x40, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y
				0x00, 0x00, 0x00, 0x04, // len(text)
				0x54, 0x65, 0x73, 0x74, // text
			},
		},
		{
			"StrokeTextMaxWidth",
			func(ctx *Context) { ctx.StrokeTextMaxWidth("Test", 32, 42, 80) },
			[]byte{
				0x3a,
				0x40, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // x
				0x40, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y
				0x40, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // maxWidth
				0x00, 0x00, 0x00, 0x04, // len(text)
				0x54, 0x65, 0x73, 0x74, // text
			},
		},
		{
			"StrokeRect",
			func(ctx *Context) { ctx.StrokeRect(93, 105, 60, 50) },
			[]byte{
				0x32,
				0x40, 0x57, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x5a, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x4e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x49, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"Translate",
			func(ctx *Context) { ctx.Translate(14.5, 23.9) },
			[]byte{
				0x38,
				0x40, 0x2d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x37, 0xe6, 0x66, 0x66, 0x66, 0x66, 0x66,
			},
		},
		{
			"Transform",
			func(ctx *Context) { ctx.Transform(0.5, 3.2, 1.7, 5.2, 11, 21) },
			[]byte{
				0x37,
				0x3f, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x09, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a,
				0x3f, 0xfb, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
				0x40, 0x14, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xcd,
				0x40, 0x26, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x35, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"SetTransform",
			func(ctx *Context) { ctx.SetTransform(0.5, 3.2, 1.7, 5.2, 11, 21) },
			[]byte{
				0x2c,
				0x3f, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x09, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a,
				0x3f, 0xfb, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
				0x40, 0x14, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xcd,
				0x40, 0x26, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x35, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"SetLineDash",
			func(ctx *Context) {
				ctx.SetLineDash([]float64{
					0.5, 1.5, 0.25, 2,
				})
			},
			[]byte{
				0x2b,
				0x00, 0x00, 0x00, 0x04, // len(segments)
				0x3f, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x3f, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x3f, 0xd0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"PutImageData",
			func(ctx *Context) {
				img := &ImageData{id: 255, ctx: ctx}
				ctx.PutImageData(img, 30, 45)
			},
			[]byte{
				0x24,                   // PutImageData
				0x00, 0x00, 0x00, 0xff, // ID
				0x40, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dx
				0x40, 0x46, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, // dy
			},
		},
		{
			"PutImageDataDirty",
			func(ctx *Context) {
				img := &ImageData{id: 256, ctx: ctx}
				ctx.PutImageDataDirty(img, 320, 200, 110, 85, 50, 40)
			},
			[]byte{
				0x3e,                   // PutImageDataDirty
				0x00, 0x00, 0x01, 0x00, // ID
				0x40, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dx
				0x40, 0x69, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dy
				0x40, 0x5b, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, // dirtyX
				0x40, 0x55, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, // dirtyY
				0x40, 0x49, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dirtyWidth
				0x40, 0x44, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dirtyHeight
			},
		},
		{
			"DrawImage",
			func(ctx *Context) {
				img := &ImageData{id: 3, ctx: ctx}
				ctx.DrawImage(img, 80, 90)
			},
			[]byte{
				0x0d,                   // DrawImage
				0x00, 0x00, 0x00, 0x03, // ID
				0x40, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dx
				0x40, 0x56, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, // dy
			},
		},
		{
			"DrawImageScaled",
			func(ctx *Context) {
				img := &ImageData{id: 1000000, ctx: ctx}
				ctx.DrawImageScaled(img, 400, 500, 2.5, 3.2)
			},
			[]byte{
				0x3f,                   // DrawImageScaled
				0x00, 0x0f, 0x42, 0x40, // ID
				0x40, 0x79, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dx
				0x40, 0x7f, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, // dy
				0x40, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dWidth
				0x40, 0x09, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a, // dHeight
			},
		},
		{
			"DrawImageSubRectangle",
			func(ctx *Context) {
				img := &ImageData{id: 65535, ctx: ctx}
				ctx.DrawImageSubRectangle(img, 42, 32, 30, 24, 10, 15, 40, 34)
			},
			[]byte{
				0x40,                   // DrawImageSubRectangle
				0x00, 0x00, 0xff, 0xff, // ID
				0x40, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // sx
				0x40, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // sy
				0x40, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // sWidth
				0x40, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // sHeight
				0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dx
				0x40, 0x2e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dy
				0x40, 0x44, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dWidth
				0x40, 0x41, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dHeight
			},
		},
		{
			"ImageData.Release",
			func(ctx *Context) {
				img := &ImageData{id: 999, ctx: ctx}
				img.Release()
			},
			[]byte{
				0x41,                   // ReleaseImage
				0x00, 0x00, 0x03, 0xe7, // ID
			},
		},
		{
			"ImageData.Release idempotency",
			func(ctx *Context) {
				img := &ImageData{id: 2, ctx: ctx}
				img.Release()
				img.Release()
			},
			[]byte{
				0x41,                   // ReleaseImage
				0x00, 0x00, 0x00, 0x02, // ID
			},
		},
		{
			"CreateLinearGradient",
			func(ctx *Context) {
				ctx.CreateLinearGradient(70, 80, 100, 120)
			},
			[]byte{
				0x09,                   // CreateLinearGradient
				0x00, 0x00, 0x00, 0x00, // ID
				0x40, 0x51, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, // x0
				0x40, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y0
				0x40, 0x59, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // x1
				0x40, 0x5e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y1
			},
		},
		{
			"CreateLinearGradient: ID generation",
			func(ctx *Context) {
				ctx.CreateLinearGradient(0, 0, 0, 0)
				ctx.CreateLinearGradient(0, 0, 0, 0)
				ctx.CreateLinearGradient(0, 0, 0, 0)
			},
			[]byte{
				0x09,                   // CreateLinearGradient
				0x00, 0x00, 0x00, 0x00, // ID = 0
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // x0
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y0
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // x1
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y1
				0x09,                   // CreateLinearGradient
				0x00, 0x00, 0x00, 0x01, // ID = 1
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // x0
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y0
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // x1
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y1
				0x09,                   // CreateLinearGradient
				0x00, 0x00, 0x00, 0x02, // ID = 2
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // x0
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y0
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // x1
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y1
			},
		},
		{
			"CreateRadialGradient",
			func(ctx *Context) {
				ctx.CreateRadialGradient(70, 80, 30, 100, 120, 20)
			},
			[]byte{
				0x0b,                   // CreateRadialGradient
				0x00, 0x00, 0x00, 0x00, // ID
				0x40, 0x51, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, // x0
				0x40, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y0
				0x40, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // r0
				0x40, 0x59, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // x1
				0x40, 0x5e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // y1
				0x40, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // r1
			},
		},
		{
			"Gradient.AddColorStop",
			func(ctx *Context) {
				g := &Gradient{id: 100000, ctx: ctx}
				g.AddColorStop(4.6, color.RGBA{R: 25, G: 110, B: 48, A: 255})
			},
			[]byte{
				0x14,                   // GradientAddColorStop
				0x00, 0x01, 0x86, 0xa0, // ID
				0x40, 0x12, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, // offset
				0x19, 0x6e, 0x30, 0xff, // R, G, B, A
			},
		},
		{
			"Gradient.AddColorStopString",
			func(ctx *Context) {
				g := &Gradient{id: 12345, ctx: ctx}
				g.AddColorStopString(23.4, "yellow")
			},
			[]byte{
				0x15,                   // GradientAddColorStopString
				0x00, 0x00, 0x30, 0x39, // ID
				0x40, 0x37, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
				0x00, 0x00, 0x00, 0x06, // len(color)
				0x79, 0x65, 0x6c, 0x6c, 0x6f, 0x77, // color
			},
		},
		{
			"Gradient.Release",
			func(ctx *Context) {
				g := &Gradient{id: 132, ctx: ctx}
				g.Release()
			},
			[]byte{
				0x21,                   // ReleaseGradient
				0x00, 0x00, 0x00, 0x84, // ID
			},
		},
		{
			"Gradient.Release idempotency",
			func(ctx *Context) {
				g := &Gradient{id: 4, ctx: ctx}
				g.Release()
				g.Release()
			},
			[]byte{
				0x21,                   // ReleaseGradient
				0x00, 0x00, 0x00, 0x04, // ID
			},
		},
		{
			"CreatePattern",
			func(ctx *Context) {
				img := &ImageData{id: 16}
				ctx.CreatePattern(img, PatternRepeatX)
				ctx.CreatePattern(img, PatternRepeatY)
				ctx.CreatePattern(img, PatternNoRepeat)
			},
			[]byte{
				0x0a,                   // CreatePattern
				0x00, 0x00, 0x00, 0x00, // ID
				0x00, 0x00, 0x00, 0x10, // ImageData ID
				0x01,
				0x0a,                   // CreatePattern
				0x00, 0x00, 0x00, 0x01, // ID
				0x00, 0x00, 0x00, 0x10, // ImageData ID
				0x02,
				0x0a,                   // CreatePattern
				0x00, 0x00, 0x00, 0x02, // ID
				0x00, 0x00, 0x00, 0x10, // ImageData ID
				0x03,
			},
		},
		{
			"Pattern.Release",
			func(ctx *Context) {
				p := &Pattern{id: 1023, ctx: ctx}
				p.Release()
			},
			[]byte{
				0x1b,                   // ReleasePattern
				0x00, 0x00, 0x03, 0xff, // ID
			},
		},
		{
			"Pattern.Release idempotency",
			func(ctx *Context) {
				p := &Pattern{id: 10, ctx: ctx}
				p.Release()
				p.Release()
			},
			[]byte{
				0x1b,                   // ReleasePattern
				0x00, 0x00, 0x00, 0x0a, // ID
			},
		},
		{
			"GetImageData",
			func(ctx *Context) {
				ctx.GetImageData(10, 10, 320, 200)
			},
			[]byte{
				0x44,                   // GetImageData
				0x00, 0x00, 0x00, 0x00, // ID
				0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // sx
				0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // sy
				0x40, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // sw
				0x40, 0x69, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // sh
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			draws := make(chan []byte)
			ctx := newContext(draws, nil, config{})
			go func(draw func(*Context)) {
				draw(ctx)
				ctx.Flush()
			}(tt.draw)
			got := <-draws
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot : %#02v\nwant: %#02v", got, tt.want)
			}
		})
	}
}

func TestImageDataSize(t *testing.T) {
	tests := []struct {
		width  int
		height int
	}{
		{12, 32},
		{320, 200},
		{1350, 875},
	}
	ctx := newContext(nil, nil, config{})
	for _, tt := range tests {
		got := ctx.CreateImageData(image.NewRGBA(image.Rect(0, 0, tt.width, tt.height)))
		if got.Width() != tt.width || got.Height() != tt.height {
			t.Errorf("got: W %d H %d, want: W %d H %d",
				got.Width(), got.Height(),
				tt.width, tt.height)
		}
	}
}

func TestCanvasSize(t *testing.T) {
	tests := []struct {
		width, height int
	}{
		{width: 0, height: 0},
		{width: 12, height: 36},
		{width: 840, height: 900},
		{width: 42314, height: 42355},
	}
	for _, tt := range tests {
		cfg := configFrom([]Option{Size(tt.width, tt.height)})
		ctx := newContext(nil, nil, cfg)
		gotWidth := ctx.CanvasWidth()
		gotHeight := ctx.CanvasHeight()
		if gotWidth != tt.width || gotHeight != tt.height {
			t.Errorf("got: W %d H %d, want: W %d H %d",
				gotWidth, gotHeight, tt.width, tt.height)
		}
	}
}

func TestUseAfterRelease(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		draw     func(ctx *Context)
	}{
		{
			name:     "ImageData",
			typeName: "ImageData",
			draw: func(ctx *Context) {
				img := ctx.CreateImageData(image.NewRGBA(image.Rect(0, 0, 0, 0)))
				img.Release()
				ctx.DrawImage(img, 0, 0)
			},
		},
		{
			name:     "Gradient",
			typeName: "Gradient",
			draw: func(ctx *Context) {
				g := ctx.CreateLinearGradient(0, 0, 0, 0)
				g.Release()
				g.AddColorStop(0, color.Black)
			},
		},
		{
			name:     "Pattern",
			typeName: "Pattern",
			draw: func(ctx *Context) {
				img := ctx.CreateImageData(image.NewRGBA(image.Rect(0, 0, 0, 0)))
				p := ctx.CreatePattern(img, PatternRepeat)
				p.Release()
				ctx.SetFillStylePattern(p)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("expected panic, but did not panic")
					return
				}
				want := tt.typeName + ": use after release"
				if r != want {
					t.Errorf("expected panic message %q, but was: %q", want, r)
				}
			}()
			ctx := newContext(nil, nil, config{})
			tt.draw(ctx)
		})
	}
}
