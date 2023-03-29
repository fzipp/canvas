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
	"image/color"
)

// Context is the server-side drawing context for a client-side canvas. It
// buffers all drawing operations until the Flush method is called. The Flush
// method then sends the buffered operations to the client.
//
// The Context for a client-server connection is obtained from the parameter
// of the run function that was passed to ListenAndServe, ListenAndServeTLS,
// or NewServeMux.
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

// Events returns a channel of events sent by the client.
//
// A type switch on the received Event values can differentiate between the
// concrete event types such as MouseDownEvent or KeyUpEvent.
func (ctx *Context) Events() <-chan Event {
	return ctx.events
}

// CanvasWidth returns the width of the canvas in pixels.
func (ctx *Context) CanvasWidth() int {
	return ctx.config.width
}

// CanvasHeight returns the height of the canvas in pixels.
func (ctx *Context) CanvasHeight() int {
	return ctx.config.height
}

// SetFillStyle sets the color to use inside shapes.
// The default color is black.
func (ctx *Context) SetFillStyle(c color.Color) {
	ctx.buf.addByte(bFillStyle)
	ctx.buf.addColor(c)
}

// SetFillStyleString sets the color to use inside shapes.
// The color is parsed as a CSS color value like "#a100cb", "#ccc",
// "darkgreen", "rgba(0.5, 0.2, 0.7, 1.0)", etc.
// The default color is "#000" (black).
func (ctx *Context) SetFillStyleString(color string) {
	ctx.buf.addByte(bFillStyleString)
	ctx.buf.addString(color)
}

// SetFillStyleGradient sets the gradient (a linear or radial gradient) to
// use inside shapes.
func (ctx *Context) SetFillStyleGradient(g *Gradient) {
	g.checkUseAfterRelease()
	ctx.buf.addByte(bFillStyleGradient)
	ctx.buf.addUint32(g.id)
}

// SetFillStylePattern sets the pattern (a repeating image) to use inside
// shapes.
func (ctx *Context) SetFillStylePattern(p *Pattern) {
	p.checkUseAfterRelease()
	ctx.buf.addByte(bFillStylePattern)
	ctx.buf.addUint32(p.id)
}

// SetFont sets the current text style to use when drawing text. This string
// uses the same syntax as the CSS font specifier. The default font is
// "10px sans-serif".
func (ctx *Context) SetFont(font string) {
	ctx.buf.addByte(bFont)
	ctx.buf.addString(font)
}

// SetGlobalAlpha sets the alpha (transparency) value that is applied to
// shapes and images before they are drawn onto the canvas.
// The alpha value is a number between 0.0 (fully transparent) and 1.0
// (fully opaque), inclusive. The default value is 1.0.
// Values outside that range, including ±Inf and NaN, will not be set.
func (ctx *Context) SetGlobalAlpha(alpha float64) {
	ctx.buf.addByte(bGlobalAlpha)
	ctx.buf.addFloat64(alpha)
}

// SetGlobalCompositeOperation sets the type of compositing operation to
// apply when drawing new shapes.
//
// The default mode is OpSourceOver.
func (ctx *Context) SetGlobalCompositeOperation(mode CompositeOperation) {
	ctx.buf.addByte(bGlobalCompositeOperation)
	ctx.buf.addByte(byte(mode))
}

// SetImageSmoothingEnabled determines whether scaled images are smoothed
// (true, default) or not (false).
//
// This property is useful for games and other apps that use pixel art. When
// enlarging images, the default resizing algorithm will blur the pixels. Set
// this property to false to retain the pixels' sharpness.
func (ctx *Context) SetImageSmoothingEnabled(enabled bool) {
	ctx.buf.addByte(bImageSmoothingEnabled)
	ctx.buf.addBool(enabled)
}

// SetLineCap sets the shape used to draw the end points of lines.
// The default value is CapButt.
//
// Note: Lines can be drawn with the Stroke, StrokeRect, and StrokeText
// methods.
func (ctx *Context) SetLineCap(cap LineCap) {
	ctx.buf.addByte(bLineCap)
	ctx.buf.addByte(byte(cap))
}

// SetLineDashOffset sets the line dash offset, or "phase."
// The default value is 0.0.
//
// Note: Lines are drawn by calling the Stroke method.
func (ctx *Context) SetLineDashOffset(offset float64) {
	ctx.buf.addByte(bLineDashOffset)
	ctx.buf.addFloat64(offset)
}

// SetLineJoin sets the shape used to join two line segments where they meet.
// The default is JoinMiter.
//
// This property has no effect wherever two connected segments have the same
// direction, because no joining area will be added in this case. Degenerate
// segments with a length of zero (i.e., with all endpoints and control points
// at the exact same position) are also ignored.
//
// Note: Lines can be drawn with the Stroke, StrokeRect, and StrokeText
// methods.
func (ctx *Context) SetLineJoin(join LineJoin) {
	ctx.buf.addByte(bLineJoin)
	ctx.buf.addByte(byte(join))
}

// SetLineWidth sets the thickness of lines.
// The width is a number in coordinate space units.
// Zero, negative, ±Inf, and NaN values are ignored.
// This value is 1.0 by default.
//
// Note: Lines can be drawn with the Stroke, StrokeRect, and StrokeText
// methods.
func (ctx *Context) SetLineWidth(width float64) {
	ctx.buf.addByte(bLineWidth)
	ctx.buf.addFloat64(width)
}

// SetMiterLimit sets the miter limit ratio.
// The miter limit ratio is a number in coordinate space units.
// Zero, negative, ±Inf, and NaN values are ignored.
// The default value is 10.0.
func (ctx *Context) SetMiterLimit(value float64) {
	ctx.buf.addByte(bMiterLimit)
	ctx.buf.addFloat64(value)
}

// SetShadowBlur sets the amount of blur applied to shadows.
// The default is 0 (no blur).
//
// The blur level is a non-negative float specifying the level of shadow blur,
// where 0 represents no blur and larger numbers represent increasingly more
// blur. This value doesn't correspond to a number of pixels, and is not
// affected by the current transformation matrix.
// Negative, ±Inf, and NaN values are ignored.
//
// Note: Shadows are only drawn if the SetShadowColor / SetShadowColorString
// property is set to a non-transparent value. One of the SetShadowBlur,
// SetShadowOffsetX, or SetShadowOffsetY properties must be non-zero, as well.
func (ctx *Context) SetShadowBlur(level float64) {
	ctx.buf.addByte(bShadowBlur)
	ctx.buf.addFloat64(level)
}

// SetShadowColor sets the color of shadows.
// The default value is fully-transparent black.
//
// Be aware that the shadow's rendered opacity will be affected by the opacity
// of the SetFillStyle color when filling, and of the SetStrokeStyle color
// when stroking.
//
// Note: Shadows are only drawn if the SetShadowColor / SetShadowColorString
// property is set to a non-transparent value. One of the SetShadowBlur,
// SetShadowOffsetX, or SetShadowOffsetY properties must be non-zero, as well.
func (ctx *Context) SetShadowColor(c color.Color) {
	ctx.buf.addByte(bShadowColor)
	ctx.buf.addColor(c)
}

// SetShadowColorString sets the color of shadows.
// The default value is fully-transparent black.
//
// The color is parsed as a CSS color value like "#a100cb", "#ccc",
// "darkgreen", "rgba(0.5, 0.2, 0.7, 1.0)", etc.
//
// Be aware that the shadow's rendered opacity will be affected by the opacity
// of the SetFillStyle color when filling, and of the SetStrokeStyle color
// when stroking.
//
// Note: Shadows are only drawn if the SetShadowColor / SetShadowColorString
// property is set to a non-transparent value. One of the SetShadowBlur,
// SetShadowOffsetX, or SetShadowOffsetY properties must be non-zero, as well.
func (ctx *Context) SetShadowColorString(color string) {
	ctx.buf.addByte(bShadowColorString)
	ctx.buf.addString(color)
}

// SetShadowOffsetX sets the distance that shadows will be offset horizontally.
//
// The offset is a float specifying the distance that shadows will be offset
// horizontally. Positive values are to the right, and negative to the left.
// The default value is 0 (no horizontal offset). ±Inf and NaN values are
// ignored.
//
// Note: Shadows are only drawn if the SetShadowColor / SetShadowColorString
// property is set to a non-transparent value. One of the SetShadowBlur,
// SetShadowOffsetX, or SetShadowOffsetY properties must be non-zero, as well.
func (ctx *Context) SetShadowOffsetX(offset float64) {
	ctx.buf.addByte(bShadowOffsetX)
	ctx.buf.addFloat64(offset)
}

// SetShadowOffsetY sets the distance that shadows will be offset vertically.
//
// The offset is a float specifying the distance that shadows will be offset
// vertically. Positive values are down, and negative are up.
// The default value is 0 (no vertical offset). ±Inf and NaN values are
// ignored.
//
// Note: Shadows are only drawn if the SetShadowColor / SetShadowColorString
// property is set to a non-transparent value. One of the SetShadowBlur,
// SetShadowOffsetX, or SetShadowOffsetY properties must be non-zero, as well.
func (ctx *Context) SetShadowOffsetY(offset float64) {
	ctx.buf.addByte(bShadowOffsetY)
	ctx.buf.addFloat64(offset)
}

// SetStrokeStyle sets the color to use for the strokes (outlines) around
// shapes. The default color is black.
//
// The color is parsed as a CSS color value like "#a100cb", "#ccc",
// "darkgreen", "rgba(0.5, 0.2, 0.7, 1.0)", etc.
func (ctx *Context) SetStrokeStyle(c color.Color) {
	ctx.buf.addByte(bStrokeStyle)
	ctx.buf.addColor(c)
}

// SetStrokeStyleString sets the color to use for the strokes (outlines) around
// shapes. The default color is black.
func (ctx *Context) SetStrokeStyleString(color string) {
	ctx.buf.addByte(bStrokeStyleString)
	ctx.buf.addString(color)
}

// SetStrokeStyleGradient sets the gradient (a linear or radial gradient) to
// use for the strokes (outlines) around shapes.
func (ctx *Context) SetStrokeStyleGradient(g *Gradient) {
	g.checkUseAfterRelease()
	ctx.buf.addByte(bStrokeStyleGradient)
	ctx.buf.addUint32(g.id)
}

// SetStrokeStylePattern sets the pattern (a repeating image) to use for the
// strokes (outlines) around shapes.
func (ctx *Context) SetStrokeStylePattern(p *Pattern) {
	p.checkUseAfterRelease()
	ctx.buf.addByte(bStrokeStylePattern)
	ctx.buf.addUint32(p.id)
}

// SetTextAlign sets the current text alignment used when drawing text.
//
// The alignment is relative to the x value of the FillText method. For
// example, if the text alignment is set to AlignCenter, then the text's left
// edge will be at x - (textWidth / 2).
//
// The default value is AlignStart.
func (ctx *Context) SetTextAlign(align TextAlign) {
	ctx.buf.addByte(bTextAlign)
	ctx.buf.addByte(byte(align))
}

// SetTextBaseline sets the current text baseline used when drawing text.
//
// The default value is BaselineAlphabetic.
func (ctx *Context) SetTextBaseline(baseline TextBaseline) {
	ctx.buf.addByte(bTextBaseline)
	ctx.buf.addByte(byte(baseline))
}

// Arc adds a circular arc to the current sub-path.
//
// It creates a circular arc centered at (x, y) with a radius of radius,
// which must be positive. The path starts at startAngle, ends at endAngle,
// and travels in the direction given by anticlockwise. The angles are
// in radians, measured from the positive x-axis.
func (ctx *Context) Arc(x, y, radius, startAngle, endAngle float64, anticlockwise bool) {
	ctx.buf.addByte(bArc)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(radius)
	ctx.buf.addFloat64(startAngle)
	ctx.buf.addFloat64(endAngle)
	ctx.buf.addBool(anticlockwise)
}

// ArcTo adds a circular arc to the current sub-path, using the given control
// points and radius. The arc is automatically connected to the path's latest
// point with a straight line, if necessary for the specified parameters.
//
// This method is commonly used for making rounded corners.
//
// (x1, y1) are the coordinates of the first control point,
// (x2, y2) are the coordinates of the second control point.
// The radius must be non-negative.
//
// Note: Be aware that you may get unexpected results when using a relatively
// large radius: the arc's connecting line will go in whatever direction it
// must to meet the specified radius.
func (ctx *Context) ArcTo(x1, y1, x2, y2, radius float64) {
	ctx.buf.addByte(bArcTo)
	ctx.buf.addFloat64(x1)
	ctx.buf.addFloat64(y1)
	ctx.buf.addFloat64(x2)
	ctx.buf.addFloat64(y2)
	ctx.buf.addFloat64(radius)
}

// BeginPath starts a new path by emptying the list of sub-paths. Call this
// method when you want to create a new path.
//
// Note: To create a new sub-path, i.e., one matching the current canvas state,
// you can use MoveTo.
func (ctx *Context) BeginPath() {
	ctx.buf.addByte(bBeginPath)
}

// BezierCurveTo adds a cubic Bézier curve to the current sub-path. It requires
// three points: the first two are control points and the third one is the end
// point. The starting point is the latest point in the current path, which can
// be changed using MoveTo before creating the Bézier curve.
//
// (cp1x, cp1y) are the coordinates of the first control point, (cp2x, cp2y)
// the coordinates of the second control point, and (x, y) the coordinates
// of the end point.
func (ctx *Context) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	ctx.buf.addByte(bBezierCurveTo)
	ctx.buf.addFloat64(cp1x)
	ctx.buf.addFloat64(cp1y)
	ctx.buf.addFloat64(cp2x)
	ctx.buf.addFloat64(cp2y)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

// ClearRect erases the pixels in a rectangular area by setting them to
// transparent black.
//
// It sets the pixels in a rectangular area to transparent black
// (rgba(0,0,0,0)). The rectangle's corner is at (x, y), and its size is
// specified by width and height.
//
// Note: Be aware that ClearRect may cause unintended side effects if you're
// not using paths properly. Make sure to call BeginPath before starting to
// draw new items after calling ClearRect.
func (ctx *Context) ClearRect(x, y, width, height float64) {
	ctx.buf.addByte(bClearRect)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(width)
	ctx.buf.addFloat64(height)
}

// Clip turns the current or given path into the current clipping region.
// The previous clipping region, if any, is intersected with the current or
// given path to create the new clipping region.
//
// Note: Be aware that the clipping region is only constructed from shapes
// added to the path. It doesn't work with shape primitives drawn directly to
// the canvas, such as FillRect. Instead, you'd have to use Rect to add a
// rectangular shape to the path before calling Clip.
func (ctx *Context) Clip() {
	ctx.buf.addByte(bClip)
}

// ClosePath attempts to add a straight line from the current point to the
// start of the current sub-path. If the shape has already been closed or has
// only one point, this function does nothing.
//
// This method doesn't draw anything to the canvas directly. You can render the
// path using the Stroke or Fill methods.
func (ctx *Context) ClosePath() {
	ctx.buf.addByte(bClosePath)
}

// Ellipse adds an elliptical arc to the current sub-path.
//
// It creates an elliptical arc centered at (x, y) with the radii radiusX and
// radiusY. The path starts at startAngle and ends at endAngle, and travels in
// the direction given by anticlockwise.
//
// The radii must be non-negative. The rotation and the angles are expressed
// in radians. The angles are measured clockwise from the positive x-axis.
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

// Fill fills the current path with the current fill style (see SetFillStyle,
// SetFillStyleString, SetFillStyleGradient, SetFillStylePattern).
func (ctx *Context) Fill() {
	ctx.buf.addByte(bFill)
}

// FillRect draws a rectangle that is filled according to the current fill
// style.
//
// It draws a filled rectangle whose starting point is at (x, y) and whose
// size is specified by width and height. The fill style is determined by the
// current fill style (see SetFillStyle, SetFillStyleString,
// SetFillStyleGradient, SetFillStylePattern).
//
// This method draws directly to the canvas without modifying the current path,
// so any subsequent Fill or Stroke calls will have no effect on it.
func (ctx *Context) FillRect(x, y, width, height float64) {
	ctx.buf.addByte(bFillRect)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(width)
	ctx.buf.addFloat64(height)
}

// FillText draws a text string at the specified coordinates, filling the
// string's characters with the current fill style (see SetFillStyle,
// SetFillStyleString, SetFillStyleGradient, SetFillStylePattern).
//
// This method draws directly to the canvas without modifying the current path,
// so any subsequent Fill or Stroke calls will have no effect on it.
//
// The text is rendered using the font and text layout configuration as defined
// by the SetFont, SetTextAlign, and SetTextBaseline properties.
func (ctx *Context) FillText(text string, x, y float64) {
	ctx.buf.addByte(bFillText)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addString(text)
}

// FillTextMaxWidth draws a text string at the specified coordinates, filling
// the string's characters with the current fill style (see SetFillStyle,
// SetFillStyleString, SetFillStyleGradient, SetFillStylePattern).
//
// The maxWidth parameter specifies the maximum number of pixels wide the text
// may be once rendered. The user agent will adjust the kerning, select a more
// horizontally condensed font (if one is available or can be generated without
// loss of quality), or scale down to a smaller font size in order to fit the
// text in the specified width.
//
// This method draws directly to the canvas without modifying the current path,
// so any subsequent Fill or Stroke calls will have no effect on it.
//
// The text is rendered using the font and text layout configuration as defined
// by the SetFont, SetTextAlign, and SetTextBaseline properties.
func (ctx *Context) FillTextMaxWidth(text string, x, y, maxWidth float64) {
	ctx.buf.addByte(bFillTextMaxWidth)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(maxWidth)
	ctx.buf.addString(text)
}

// LineTo adds a straight line to the current sub-path by connecting the
// sub-path's last point to the specified (x, y) coordinates.
//
// Like other methods that modify the current path, this method does not
// directly render anything. To draw the path onto a canvas, you can use the
// Fill or Stroke methods.
func (ctx *Context) LineTo(x, y float64) {
	ctx.buf.addByte(bLineTo)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

// MoveTo begins a new sub-path at the point specified by the given (x, y)
// coordinates.
func (ctx *Context) MoveTo(x, y float64) {
	ctx.buf.addByte(bMoveTo)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

// QuadraticCurveTo adds a quadratic Bézier curve to the current sub-path.
// It requires two points: the first one is a control point and the second one
// is the end point. The starting point is the latest point in the current
// path, which can be changed using MoveTo before creating the quadratic
// Bézier curve.
//
// (cpx, cpy) is the coordinate of the control point, and (x, y) is the
// coordinate of the end point.
func (ctx *Context) QuadraticCurveTo(cpx, cpy, x, y float64) {
	ctx.buf.addByte(bQuadraticCurveTo)
	ctx.buf.addFloat64(cpx)
	ctx.buf.addFloat64(cpy)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

// Rect adds a rectangle to the current path.
//
// It creates a rectangular path whose starting point is at (x, y) and whose
// size is specified by width and height.
// Like other methods that modify the current path, this method does not
// directly render anything. To draw the rectangle onto a canvas, you can use
// the Fill or Stroke methods.
//
// Note: To both create and render a rectangle in one step, use the FillRect
// or StrokeRect methods.
func (ctx *Context) Rect(x, y, width, height float64) {
	ctx.buf.addByte(bRect)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(width)
	ctx.buf.addFloat64(height)
}

// Restore restores the most recently saved canvas state by popping the top
// entry in the drawing state stack. If there is no saved state, this method
// does nothing.
//
// For more information about the drawing state, see Save.
func (ctx *Context) Restore() {
	ctx.buf.addByte(bRestore)
}

// Rotate adds a rotation to the transformation matrix.
//
// The rotation angle is expressed in radians (clockwise). You can use
// degree * math.Pi / 180 to calculate a radian from a degree.
//
// The rotation center point is always the canvas origin. To change the center
// point, you will need to move the canvas by using the Translate method.
func (ctx *Context) Rotate(angle float64) {
	ctx.buf.addByte(bRotate)
	ctx.buf.addFloat64(angle)
}

// Save saves the entire state of the canvas by pushing the current state onto
// a stack.
//
// The drawing state that gets saved onto a stack consists of:
//   - The current transformation matrix.
//   - The current clipping region.
//   - The current dash list.
//   - The current values of the attributes set via the following methods:
//     SetStrokeStyle*, SetFillStyle*, SetGlobalAlpha, SetLineWidth,
//     SetLineCap, SetLineJoin, SetMiterLimit, SetLineDashOffset,
//     SetShadowOffsetX, SetShadowOffsetY, SetShadowBlur, SetShadowColor*,
//     SetGlobalCompositeOperation, SetFont, SetTextAlign, SetTextBaseline,
//     SetImageSmoothingEnabled.
func (ctx *Context) Save() {
	ctx.buf.addByte(bSave)
}

// Scale adds a scaling transformation to the canvas units horizontally
// and/or vertically.
//
// By default, one unit on the canvas is exactly one pixel. A scaling
// transformation modifies this behavior. For instance, a scaling factor of 0.5
// results in a unit size of 0.5 pixels; shapes are thus drawn at half the
// normal size. Similarly, a scaling factor of 2.0 increases the unit size so
// that one unit becomes two pixels; shapes are thus drawn at twice the normal
// size.
//
// x is the scaling factor in the horizontal direction. A negative value flips
// pixels across the vertical axis. A value of 1 results in no horizontal
// scaling.
//
// y is the scaling factor in the vertical direction. A negative value flips
// pixels across the horizontal axis. A value of 1 results in no vertical
// scaling.
func (ctx *Context) Scale(x, y float64) {
	ctx.buf.addByte(bScale)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

// Stroke outlines the current or given path with the current stroke style.
//
// Strokes are aligned to the center of a path; in other words, half of the
// stroke is drawn on the inner side, and half on the outer side.
//
// The stroke is drawn using the [non-zero winding rule], which means that path
// intersections will still get filled.
//
// [non-zero winding rule]: https://en.wikipedia.org/wiki/Nonzero-rule
func (ctx *Context) Stroke() {
	ctx.buf.addByte(bStroke)
}

// StrokeText strokes - that is, draws the outlines of - the characters of a
// text string at the specified coordinates.
//
// The text is rendered using the settings specified by SetFont, SetTextAlign,
// and SetTextBaseline. (x, y) is the coordinate of the point at which to begin
// drawing the text.
//
// This method draws directly to the canvas without modifying the current path,
// so any subsequent Fill or Stroke calls will have no effect on it.
//
// Use the FillText method to fill the text characters rather than having just
// their outlines drawn.
func (ctx *Context) StrokeText(text string, x, y float64) {
	ctx.buf.addByte(bStrokeText)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addString(text)
}

// StrokeTextMaxWidth strokes - that is, draws the outlines of - the characters
// of a text string at the specified coordinates. A parameter allows
// specifying a maximum width for the rendered text, which the user agent will
// achieve by condensing the text or by using a lower font size.
//
// The text is rendered using the settings specified by SetFont, SetTextAlign,
// and SetTextBaseline. (x, y) is the coordinate of the point at which to begin
// drawing the text.
//
// The user agent will adjust the kerning, select a more horizontally condensed
// font (if one is available or can be generated without loss of quality), or
// scale down to a smaller font size in order to fit the text in the specified
// maxWidth.
//
// This method draws directly to the canvas without modifying the current path,
// so any subsequent Fill or Stroke calls will have no effect on it.
//
// Use the FillText method to fill the text characters rather than having just
// their outlines drawn.
func (ctx *Context) StrokeTextMaxWidth(text string, x, y, maxWidth float64) {
	ctx.buf.addByte(bStrokeTextMaxWidth)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(maxWidth)
	ctx.buf.addString(text)
}

// StrokeRect draws a rectangle that is stroked (outlined) according to the
// current stroke style and other context settings.
//
// It draws a stroked rectangle whose starting point is at (x, y) and whose
// size is specified by width and height.
//
// This method draws directly to the canvas without modifying the current path
// so any subsequent Fill or Stroke calls will have no effect on it.
func (ctx *Context) StrokeRect(x, y, width, height float64) {
	ctx.buf.addByte(bStrokeRect)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
	ctx.buf.addFloat64(width)
	ctx.buf.addFloat64(height)
}

// Translate adds a translation transformation to the current matrix by moving
// the canvas and its origin x units horizontally and y units vertically on the
// grid.
//
// Positive x values are to the right, and negative to the left.
// Positive y values are down, and negative are up.
func (ctx *Context) Translate(x, y float64) {
	ctx.buf.addByte(bTranslate)
	ctx.buf.addFloat64(x)
	ctx.buf.addFloat64(y)
}

// Transform multiplies the current transformation with the matrix described
// by the arguments of this method. This lets you scale, rotate, translate
// (move), and skew the context. The transformation matrix is described by:
//
//	[ a b 0 ]
//	[ c d 0 ]
//	[ e f 1 ]
//
//	a: Horizontal scaling. A value of 1 results in no scaling.
//	b: Vertical skewing.
//	c: Horizontal skewing.
//	d: Vertical scaling. A value of 1 results in no scaling.
//	e: Horizontal translation (moving).
//	f: Vertical translation (moving).
//
// Note: See also the SetTransform method, which resets the current transform
// to the identity matrix and then invokes Transform.
func (ctx *Context) Transform(a, b, c, d, e, f float64) {
	ctx.buf.addByte(bTransform)
	ctx.buf.addFloat64(a)
	ctx.buf.addFloat64(b)
	ctx.buf.addFloat64(c)
	ctx.buf.addFloat64(d)
	ctx.buf.addFloat64(e)
	ctx.buf.addFloat64(f)
}

// SetTransform resets (overrides) the current transformation to the identity
// matrix, and then invokes a transformation described by the arguments of this
// method. This lets you scale, rotate, translate (move), and skew the context.
// The transformation matrix is described by:
//
//	[ a b 0 ]
//	[ c d 0 ]
//	[ e f 1 ]
//
//	a: Horizontal scaling. A value of 1 results in no scaling.
//	b: Vertical skewing.
//	c: Horizontal skewing.
//	d: Vertical scaling. A value of 1 results in no scaling.
//	e: Horizontal translation (moving).
//	f: Vertical translation (moving).
//
// Note: See also the Transform method; instead of overriding the current
// transform matrix, it multiplies it with a given one.
func (ctx *Context) SetTransform(a, b, c, d, e, f float64) {
	ctx.buf.addByte(bSetTransform)
	ctx.buf.addFloat64(a)
	ctx.buf.addFloat64(b)
	ctx.buf.addFloat64(c)
	ctx.buf.addFloat64(d)
	ctx.buf.addFloat64(e)
	ctx.buf.addFloat64(f)
}

// SetLineDash sets the line dash pattern used when stroking lines. It uses a
// slice of values that specify alternating lengths of lines and gaps which
// describe the pattern.
//
// The segments are a slice of numbers that specify distances to alternately
// draw a line and a gap (in coordinate space units). If the number of elements
// in the slice is odd, the elements of the slice get copied and concatenated.
// For example, {5, 15, 25} will become {5, 15, 25, 5, 15, 25}. If the slice
// is empty, the line dash list is cleared and line strokes return to being
// solid.
func (ctx *Context) SetLineDash(segments []float64) {
	ctx.buf.addByte(bSetLineDash)
	ctx.buf.addUint32(uint32(len(segments)))
	for _, seg := range segments {
		ctx.buf.addFloat64(seg)
	}
}

// CreateImageData creates a new, blank ImageData object on the client with the
// specified dimensions. All of the pixels in the new object are transparent
// black. The ImageData object should be released with the ImageData.Release
// method when it is no longer needed.
func (ctx *Context) CreateImageData(m image.Image) *ImageData {
	rgba := ensureRGBA(m)
	bounds := m.Bounds()
	id := ctx.imageDataIDs.generateID()
	ctx.buf.addByte(bCreateImageData)
	ctx.buf.addUint32(id)
	ctx.buf.addUint32(uint32(bounds.Dx()))
	ctx.buf.addUint32(uint32(bounds.Dy()))
	ctx.buf.addBytes(rgba.Pix)
	return &ImageData{id: id, ctx: ctx, width: bounds.Dx(), height: bounds.Dy()}
}

// PutImageData paints data from the given ImageData object onto the
// canvas. If a dirty rectangle is provided, only the pixels from that
// rectangle are painted. This method is not affected by the canvas
// transformation matrix.
//
// (dx, dy) is the position at which to place the image data in the destination
// canvas.
//
// Note: Image data can be retrieved from a canvas using the GetImageData
// method.
func (ctx *Context) PutImageData(src *ImageData, dx, dy float64) {
	src.checkUseAfterRelease()
	ctx.buf.addByte(bPutImageData)
	ctx.buf.addUint32(src.id)
	ctx.buf.addFloat64(dx)
	ctx.buf.addFloat64(dy)
}

// PutImageDataDirty paints data from the given ImageData object onto the
// canvas. If a dirty rectangle is provided, only the pixels from that
// rectangle are painted. This method is not affected by the canvas
// transformation matrix.
//
// (dx, dy) is the position at which to place the image data in the destination
// canvas.
//
// (dirtyX, dirtyY) is the position of the top-left corner from which the image
// data will be extracted; dirtyWidth and dirtyHeight are the width and height
// of the rectangle to be painted.
//
// Note: Image data can be retrieved from a canvas using the GetImageData
// method.
func (ctx *Context) PutImageDataDirty(src *ImageData, dx, dy, dirtyX, dirtyY, dirtyWidth, dirtyHeight float64) {
	src.checkUseAfterRelease()
	ctx.buf.addByte(bPutImageDataDirty)
	ctx.buf.addUint32(src.id)
	ctx.buf.addFloat64(dx)
	ctx.buf.addFloat64(dy)
	ctx.buf.addFloat64(dirtyX)
	ctx.buf.addFloat64(dirtyY)
	ctx.buf.addFloat64(dirtyWidth)
	ctx.buf.addFloat64(dirtyHeight)
}

// DrawImage draws an image onto the canvas.
//
// (dx, dy) is the position in the destination canvas at which to place the
// top-left corner of the source image.
func (ctx *Context) DrawImage(src *ImageData, dx, dy float64) {
	src.checkUseAfterRelease()
	ctx.buf.addByte(bDrawImage)
	ctx.buf.addUint32(src.id)
	ctx.buf.addFloat64(dx)
	ctx.buf.addFloat64(dy)
}

// DrawImageScaled draws an image onto the canvas.
//
// (dx, dy) is the position in the destination canvas at which to place the
// top-left corner of the source image; dWidth and dHeight are the width to
// draw the image in the destination canvas. This allows scaling of the drawn
// image.
func (ctx *Context) DrawImageScaled(src *ImageData, dx, dy, dWidth, dHeight float64) {
	src.checkUseAfterRelease()
	ctx.buf.addByte(bDrawImageScaled)
	ctx.buf.addUint32(src.id)
	ctx.buf.addFloat64(dx)
	ctx.buf.addFloat64(dy)
	ctx.buf.addFloat64(dWidth)
	ctx.buf.addFloat64(dHeight)
}

// DrawImageSubRectangle draws an image onto the canvas.
//
// (sx, sy) is the position of the top left corner of the sub-rectangle of the
// source image to draw into the destination context; sWidth and sHeight are
// the width and height of the sub-rectangle of the source image to draw into
// the destination context.
//
// (dx, dy) is the position in the destination canvas at which to place the
// top-left corner of the source image; dWidth and dHeight are the width to
// draw the image in the destination canvas. This allows scaling of the drawn
// image.
func (ctx *Context) DrawImageSubRectangle(src *ImageData, sx, sy, sWidth, sHeight, dx, dy, dWidth, dHeight float64) {
	src.checkUseAfterRelease()
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

// CreateLinearGradient creates a gradient along the line connecting two given
// coordinates. To be applied to a shape, the gradient must first be set via
// the SetFillStyleGradient or SetStrokeStyleGradient methods.
//
// (x0, y0) defines the start point, and (x1, y1) defines the end point of the
// gradient line.
//
// Note: Gradient coordinates are global, i.e., relative to the current
// coordinate space. When applied to a shape, the coordinates are NOT relative
// to the shape's coordinates.
func (ctx *Context) CreateLinearGradient(x0, y0, x1, y1 float64) *Gradient {
	id := ctx.gradientIDs.generateID()
	ctx.buf.addByte(bCreateLinearGradient)
	ctx.buf.addUint32(id)
	ctx.buf.addFloat64(x0)
	ctx.buf.addFloat64(y0)
	ctx.buf.addFloat64(x1)
	ctx.buf.addFloat64(y1)
	return &Gradient{id: id, ctx: ctx}
}

// CreateRadialGradient creates a radial gradient using the size and
// coordinates of two circles. To be applied to a shape, the gradient must
// first be set via the SetFillStyleGradient or SetStrokeStyleGradient methods.
//
// (x0, y0) defines the center, and r0 the radius of the start circle.
// (x1, y1) defines the center, and r1 the radius of the end circle.
// Each radius must be non-negative and finite.
//
// Note: Gradient coordinates are global, i.e., relative to the current
// coordinate space. When applied to a shape, the coordinates are NOT relative
// to the shape's coordinates.
func (ctx *Context) CreateRadialGradient(x0, y0, r0, x1, y1, r1 float64) *Gradient {
	id := ctx.gradientIDs.generateID()
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

// CreatePattern creates a pattern using the specified image and repetition.
// The repetition indicates how to repeat the pattern's image.
//
// This method doesn't draw anything to the canvas directly. The pattern it
// creates must be set via the SetFillStylePattern or SetStrokeStylePattern
// methods, after which it is applied to any subsequent drawing.
func (ctx *Context) CreatePattern(src *ImageData, repetition PatternRepetition) *Pattern {
	src.checkUseAfterRelease()
	id := ctx.patternIDs.generateID()
	ctx.buf.addByte(bCreatePattern)
	ctx.buf.addUint32(id)
	ctx.buf.addUint32(src.id)
	ctx.buf.addByte(byte(repetition))
	return &Pattern{id: id, ctx: ctx}
}

// GetImageData returns an ImageData object representing the underlying pixel
// data for a specified portion of the canvas.
//
// (sx, sy) is the position of the top-left corner of the rectangle from which
// the ImageData will be extracted; sw and sh are the width and height of the
// rectangle from which the ImageData will be extracted.
//
// This method is not affected by the canvas's transformation matrix. If the
// specified rectangle extends outside the bounds of the canvas, the pixels
// outside the canvas are transparent black in the returned ImageData object.
//
// Note: Image data can be painted onto a canvas using the PutImageData method.
func (ctx *Context) GetImageData(sx, sy, sw, sh float64) *ImageData {
	id := ctx.imageDataIDs.generateID()
	ctx.buf.addByte(bGetImageData)
	ctx.buf.addUint32(id)
	ctx.buf.addFloat64(sx)
	ctx.buf.addFloat64(sy)
	ctx.buf.addFloat64(sw)
	ctx.buf.addFloat64(sh)
	return &ImageData{id: id, ctx: ctx, width: int(sw), height: int(sh)}
}

// Flush sends the buffered drawing operations of the context from the server
// to the client.
//
// Nothing is displayed on the client canvas until Flush is called.
// An animation loop usually has one flush per animation frame.
func (ctx *Context) Flush() {
	ctx.draws <- ctx.buf.bytes
	ctx.buf.reset()
}

type idGenerator struct {
	next uint32
}

func (g *idGenerator) generateID() uint32 {
	id := g.next
	g.next++
	return id
}
