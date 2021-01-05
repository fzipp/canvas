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
	modKeys modifierKey
}

func (e MouseEvent) mask() eventMask {
	return maskMouseMove | maskMouseUp | maskKeyDown | maskClick | maskDblClick | maskAuxClick
}

func (e *MouseEvent) AltKey() bool {
	return e.isPressed(modKeyAlt)
}

func (e *MouseEvent) ShiftKey() bool {
	return e.isPressed(modKeyShift)
}

func (e *MouseEvent) CtrlKey() bool {
	return e.isPressed(modKeyCtrl)
}

func (e *MouseEvent) MetaKey() bool {
	return e.isPressed(modKeyMeta)
}

func (e *MouseEvent) isPressed(k modifierKey) bool {
	return e.modKeys&k != 0
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
	Key     string
	modKeys modifierKey
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

type modifierKey byte

const (
	modKeyAlt modifierKey = 1 << iota
	modKeyShift
	modKeyCtrl
	modKeyMeta
)

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
