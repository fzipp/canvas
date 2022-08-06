// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

// LineCap represents the shape used to draw the end points of lines.
type LineCap byte

const (
	// CapButt squares off the ends of lines at the endpoints.
	CapButt LineCap = iota
	// CapRound rounds the ends of lines.
	CapRound
	// CapSquare squares off the ends of lines by adding a box with an equal
	// width and half the height of the line's thickness.
	CapSquare
)

// LineJoin represents the shape used to join two line segments where they
// meet.
type LineJoin byte

const (
	// JoinMiter joins connected segments by extending their outside edges to
	// connect at a single point, with the effect of filling an additional
	// lozenge-shaped area. This setting is affected by Context.SetMiterLimit.
	JoinMiter LineJoin = iota
	// JoinRound rounds off the corners of a shape by filling an additional
	// sector of disc centered at the common endpoint of connected segments.
	// The radius for these rounded corners is equal to the line width.
	JoinRound
	// JoinBevel fills an additional triangular area between the common
	// endpoint of connected segments, and the separate outside rectangular
	// corners of each segment.
	JoinBevel
)

// CompositeOperation represents the type of compositing operation to apply
// when drawing new shapes.
//
// For visual explanations of the composite operations see the [MDN docs]
// for CanvasRenderingContext2D.globalCompositeOperation.
//
// [MDN docs]: https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/globalCompositeOperation#operations
type CompositeOperation byte

const (
	// OpSourceOver draws new shapes on top of the existing canvas content.
	OpSourceOver CompositeOperation = iota
	// OpSourceIn draws the new shape only where both the new shape and the
	// destination canvas overlap. Everything else is made transparent.
	OpSourceIn
	// OpSourceOut draws the new shape where it doesn't overlap the existing
	// canvas content.
	OpSourceOut
	// OpSourceAtop draws the new shape only where it overlaps the existing
	// canvas content.
	OpSourceAtop
	// OpDestinationOver draws new shapes behind the existing canvas content.
	OpDestinationOver
	// OpDestinationIn keeps the existing canvas content where both the new
	// shape and existing canvas content overlap. Everything else is made
	// transparent.
	OpDestinationIn
	// OpDestinationOut keeps the existing content where it doesn't overlap
	// the new shape.
	OpDestinationOut
	// OpDestinationAtop keeps the existing canvas content only where it
	// overlaps the new shape. The new shape is drawn behind the canvas
	// content.
	OpDestinationAtop
	// OpLighter determines the color by adding color values where both shapes
	// overlap.
	OpLighter
	// OpCopy shows only the new shape.
	OpCopy
	// OpXOR makes shapes transparent where both overlap and draws them normal
	// everywhere else.
	OpXOR
	// OpMultiply multiplies the pixels of the top layer with the corresponding
	// pixels of the bottom layer. A darker picture is the result.
	OpMultiply
	// OpScreen inverts inverts, multiplies, and inverts the pixels again.
	// A lighter picture is the result (opposite of multiply)
	OpScreen
	// OpOverlay is a combination of OpMultiply and OpScreen. Dark parts on the
	// base layer become darker, and light parts become lighter.
	OpOverlay
	// OpDarken retains the darkest pixels of both layers.
	OpDarken
	// OpLighten retains the lightest pixels of both layers.
	OpLighten
	// OpColorDodge divides the bottom layer by the inverted top layer.
	OpColorDodge
	// OpColorBurn divides the inverted bottom layer by the top layer, and
	// then inverts the result.
	OpColorBurn
	// OpHardLight is a combination of multiply and screen like overlay, but
	// with top and bottom layer swapped.
	OpHardLight
	// OpSoftLight is a softer version of hard-light. Pure black or white does
	// not result in pure black or white.
	OpSoftLight
	// OpDifference subtracts the bottom layer from the top layer or the other
	// way round to always get a positive value.
	OpDifference
	// OpExclusion is like OpDifference, but with lower contrast.
	OpExclusion
	// OpHue preserves the luma and chroma of the bottom layer, while adopting
	// the hue of the top layer.
	OpHue
	// OpSaturation preserves the luma and hue of the bottom layer, while
	// adopting the chroma of the top layer.
	OpSaturation
	// OpColor preserves the luma of the bottom layer, while adopting the hue
	// and chroma of the top layer.
	OpColor
	// OpLuminosity preserves the hue and chroma of the bottom layer, while
	// adopting the luma of the top layer.
	OpLuminosity
)

// TextAlign represents the text alignment used when drawing text.
type TextAlign byte

const (
	// AlignStart means the text is aligned at the normal start of the line
	// (left-aligned for left-to-right locales, right-aligned for right-to-left
	// locales).
	AlignStart TextAlign = iota
	// AlignEnd means the text is aligned at the normal end of the line
	// (right-aligned for left-to-right locales, left-aligned for right-to-left
	// locales).
	AlignEnd
	// AlignLeft means the text is left-aligned.
	AlignLeft
	// AlignRight means the text is right-aligned.
	AlignRight
	// AlignCenter means the text is centered.
	AlignCenter
)

// TextBaseline represents the text baseline used when drawing text.
type TextBaseline byte

const (
	// BaselineAlphabetic means the text baseline is the normal alphabetic
	// baseline.
	BaselineAlphabetic TextBaseline = iota
	// BaselineIdeographic means the text baseline is the ideographic baseline;
	// this is the bottom of the body of the characters, if the main body of
	// characters protrudes beneath the alphabetic baseline.
	// (Used by Chinese, Japanese, and Korean scripts.)
	BaselineIdeographic
	// BaselineTop means the text baseline is the top of the em square.
	BaselineTop
	// BaselineBottom means the text baseline is the bottom of the bounding
	// box. This differs from the ideographic baseline in that the ideographic
	// baseline doesn't consider descenders.
	BaselineBottom
	// BaselineHanging means the text baseline is the hanging baseline.
	// (Used by Tibetan and other Indic scripts.)
	BaselineHanging
	// BaselineMiddle means the text baseline is the middle of the em square.
	BaselineMiddle
)

// PatternRepetition indicates how to repeat a pattern's image.
type PatternRepetition byte

const (
	// PatternRepeat repeats the image in both directions.
	PatternRepeat PatternRepetition = iota
	// PatternRepeatX repeats the image only horizontally.
	PatternRepeatX
	// PatternRepeatY repeats the image only vertically.
	PatternRepeatY
	// PatternNoRepeat repeats the image in neither direction.
	PatternNoRepeat
)

const (
	bArc byte = 1 + iota
	bArcTo
	bBeginPath
	bBezierCurveTo
	bClearRect
	bClip
	bClosePath
	bCreateImageData
	bCreateLinearGradient
	bCreatePattern
	bCreateRadialGradient
	_
	bDrawImage
	bEllipse
	bFill
	bFillRect
	bFillStyle
	bFillText
	bFont
	bGradientAddColorStop
	bGradientAddColorStopString
	bFillStyleGradient
	bGlobalAlpha
	bGlobalCompositeOperation
	bImageSmoothingEnabled
	bStrokeStyleGradient
	bReleasePattern
	bLineCap
	bLineDashOffset
	bLineJoin
	bLineTo
	bLineWidth
	bReleaseGradient
	bMiterLimit
	bMoveTo
	bPutImageData
	bQuadraticCurveTo
	bRect
	bRestore
	bRotate
	bSave
	bScale
	bSetLineDash
	bSetTransform
	bShadowBlur
	bShadowColor
	bShadowOffsetX
	bShadowOffsetY
	bStroke
	bStrokeRect
	bStrokeStyle
	bStrokeText
	bTextAlign
	bTextBaseline
	bTransform
	bTranslate
	bFillTextMaxWidth
	bStrokeTextMaxWidth
	bFillStyleString
	bStrokeStyleString
	bShadowColorString
	bPutImageDataDirty
	bDrawImageScaled
	bDrawImageSubRectangle
	bReleaseImageData
	bFillStylePattern
	bStrokeStylePattern
	bGetImageData
)
