// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

type LineCap byte

const (
	CapButt LineCap = iota
	CapRound
	CapSquare
)

type LineJoin byte

const (
	JoinMiter LineJoin = iota
	JoinRound
	JoinBevel
)

type CompositeOperation byte

const (
	OpSourceOver CompositeOperation = iota
	OpSourceIn
	OpSourceOut
	OpSourceAtop
	OpDestinationOver
	OpDestinationIn
	OpDestinationOut
	OpDestinationAtop
	OpLighter
	OpCopy
	OpXOR
	OpMultiply
	OpScreen
	OpOverlay
	OpDarken
	OpLighten
	OpColorDodge
	OpColorBurn
	OpHardLight
	OpSoftLight
	OpDifference
	OpExclusion
	OpHue
	OpSaturation
	OpColor
	OpLuminosity
)

type TextAlign byte

const (
	AlignStart TextAlign = iota
	AlignEnd
	AlignLeft
	AlignRight
	AlignCenter
)

type TextBaseline byte

const (
	BaselineAlphabetic TextBaseline = iota
	BaselineIdeographic
	BaselineTop
	BaselineBottom
	BaselineHanging
	BaselineMiddle
)

type PatternRepetition byte

const (
	PatternRepeat PatternRepetition = iota
	PatternRepeatX
	PatternRepeatY
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
	bDrawFocusIfNeeded
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
	bReleaseImage
	bFillStylePattern
	bStrokeStylePattern
)
