// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image"
	"image/color"
	"testing"
)

func BenchmarkContext(b *testing.B) {
	draws := make(chan []byte)
	go func() {
		for range draws {
		}
	}()
	ctx := newContext(draws, nil, nil)
	for range b.N {
		ctx.SetFillStyle(color.White)
		ctx.SetFillStyleString("green")
		ctx.SetFillStyleGradient(&Gradient{})
		ctx.SetFillStylePattern(&Pattern{})
		ctx.SetFont("bold 48px serif")
		ctx.SetGlobalAlpha(1)
		ctx.SetGlobalCompositeOperation(OpDestinationOut)
		ctx.SetImageSmoothingEnabled(true)
		ctx.SetLineCap(CapRound)
		ctx.SetLineDashOffset(1)
		ctx.SetLineJoin(JoinBevel)
		ctx.SetLineWidth(1)
		ctx.SetMiterLimit(1)
		ctx.SetShadowBlur(1)
		ctx.SetShadowColor(color.White)
		ctx.SetShadowColorString("#ffa0b3")
		ctx.SetShadowOffsetX(1)
		ctx.SetShadowOffsetY(1)
		ctx.SetStrokeStyle(color.White)
		ctx.SetStrokeStyleString("yellow")
		ctx.SetStrokeStyleGradient(&Gradient{})
		ctx.SetStrokeStylePattern(&Pattern{})
		ctx.SetTextAlign(AlignLeft)
		ctx.SetTextBaseline(BaselineBottom)
		ctx.Arc(1, 1, 1, 1, 1, true)
		ctx.ArcTo(1, 1, 1, 1, 1)
		ctx.BeginPath()
		ctx.BezierCurveTo(1, 1, 1, 1, 1, 1)
		ctx.ClearRect(1, 1, 1, 1)
		ctx.Clip()
		ctx.ClosePath()
		ctx.Ellipse(1, 1, 1, 1, 1, 1, 1, true)
		ctx.Fill()
		ctx.FillRect(1, 1, 1, 1)
		ctx.FillText("hello, world", 1, 1)
		ctx.FillTextMaxWidth("hello, world", 1, 1, 1)
		ctx.LineTo(1, 1)
		ctx.MoveTo(1, 1)
		ctx.QuadraticCurveTo(1, 1, 1, 1)
		ctx.Rect(1, 1, 1, 1)
		ctx.Restore()
		ctx.Rotate(1)
		ctx.Save()
		ctx.Scale(1, 1)
		ctx.Stroke()
		ctx.StrokeText("hello, world", 1, 1)
		ctx.StrokeTextMaxWidth("hello, world", 1, 1, 1)
		ctx.StrokeRect(1, 1, 1, 1)
		ctx.Translate(1, 1)
		ctx.Transform(1, 1, 1, 1, 1, 1)
		ctx.SetTransform(1, 1, 1, 1, 1, 1)
		ctx.SetLineDash([]float64{1, 1, 1})
		ctx.CreateImageData(image.NewRGBA(image.Rect(0, 0, 0, 0)))
		ctx.PutImageData(&ImageData{}, 1, 1)
		ctx.PutImageDataDirty(&ImageData{}, 1, 1, 1, 1, 1, 1)
		ctx.DrawImage(&ImageData{}, 1, 1)
		ctx.DrawImageScaled(&ImageData{}, 1, 1, 1, 1)
		ctx.DrawImageSubRectangle(&ImageData{}, 1, 1, 1, 1, 1, 1, 1, 1)
		ctx.CreateLinearGradient(1, 1, 1, 1)
		ctx.CreateRadialGradient(1, 1, 1, 1, 1, 1)
		ctx.CreatePattern(&ImageData{}, PatternRepeat)
		ctx.GetImageData(1, 1, 1, 1)
		ctx.Flush()
	}
	close(draws)
}
