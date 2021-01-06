// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

type Event interface {
	mask() eventMask
}

type MouseEvent struct {
	Buttons MouseButtons
	X, Y    int
	modifierKeys
}

func (e MouseEvent) mask() eventMask {
	return maskMouseMove | maskMouseUp | maskKeyDown | maskClick | maskDblClick | maskAuxClick
}

type MouseMoveEvent struct{ MouseEvent }

func (e MouseMoveEvent) mask() eventMask { return maskMouseMove }

type MouseDownEvent struct{ MouseEvent }

func (e MouseDownEvent) mask() eventMask { return maskMouseDown }

type MouseUpEvent struct{ MouseEvent }

func (e MouseUpEvent) mask() eventMask { return maskMouseUp }

type ClickEvent struct{ MouseEvent }

func (e ClickEvent) mask() eventMask { return maskClick }

type DblClickEvent struct{ MouseEvent }

func (e DblClickEvent) mask() eventMask { return maskDblClick }

type AuxClickEvent struct{ MouseEvent }

func (e AuxClickEvent) mask() eventMask { return maskAuxClick }

type WheelEvent struct {
	MouseEvent
	DeltaX    float64
	DeltaY    float64
	DeltaZ    float64
	DeltaMode DeltaMode
}

func (e WheelEvent) mask() eventMask {
	return maskWheel
}

type DeltaMode byte

const (
	DeltaPixel DeltaMode = iota
	DeltaLine
	DeltaPage
)

type KeyboardEvent struct {
	Key string
	modifierKeys
}

func (e KeyboardEvent) mask() eventMask {
	return maskKeyPress | maskKeyDown | maskKeyUp
}

type KeyPressEvent struct{ KeyboardEvent }

func (e KeyPressEvent) mask() eventMask { return maskKeyPress }

type KeyDownEvent struct{ KeyboardEvent }

func (e KeyDownEvent) mask() eventMask { return maskKeyDown }

type KeyUpEvent struct{ KeyboardEvent }

func (e KeyUpEvent) mask() eventMask { return maskKeyUp }

type TouchEvent struct {
	Touches        TouchList
	ChangedTouches TouchList
	TargetTouches  TouchList
	modifierKeys
}

func (e TouchEvent) mask() eventMask {
	return maskTouchStart | maskTouchMove | maskTouchEnd | maskTouchCancel
}

type TouchList []Touch

type Touch struct {
	Identifier uint32
	X          int
	Y          int
}

type TouchStartEvent struct{ TouchEvent }

func (e TouchStartEvent) mask() eventMask { return maskTouchStart }

type TouchMoveEvent struct{ TouchEvent }

func (e TouchMoveEvent) mask() eventMask { return maskTouchMove }

type TouchEndEvent struct{ TouchEvent }

func (e TouchEndEvent) mask() eventMask { return maskTouchEnd }

type TouchCancelEvent struct{ TouchEvent }

func (e TouchCancelEvent) mask() eventMask { return maskTouchCancel }

type modifierKeys byte

const (
	modKeyAlt modifierKeys = 1 << iota
	modKeyShift
	modKeyCtrl
	modKeyMeta
)

func (m modifierKeys) AltKey() bool {
	return m.isPressed(modKeyAlt)
}

func (m modifierKeys) ShiftKey() bool {
	return m.isPressed(modKeyShift)
}

func (m modifierKeys) CtrlKey() bool {
	return m.isPressed(modKeyCtrl)
}

func (m modifierKeys) MetaKey() bool {
	return m.isPressed(modKeyMeta)
}

func (m modifierKeys) isPressed(k modifierKeys) bool {
	return m&k != 0
}

type eventMask int

const (
	maskMouseMove eventMask = 1 << iota
	maskMouseDown
	maskMouseUp
	maskKeyPress
	maskKeyDown
	maskKeyUp
	maskClick
	maskDblClick
	maskAuxClick
	maskWheel
	maskTouchStart
	maskTouchMove
	maskTouchEnd
	maskTouchCancel
)

type MouseButtons int

const (
	ButtonPrimary MouseButtons = 1 << iota
	ButtonSecondary
	ButtonAuxiliary
	Button4th
	Button5th
	ButtonNone MouseButtons = 0
)
