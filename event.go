// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

type Event interface {
}

type MouseEvent struct {
	Buttons MouseButtons
	X, Y    int
	modKeys modifierKey
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
type MouseDownEvent struct{ MouseEvent }
type MouseUpEvent struct{ MouseEvent }

type KeyboardEvent struct {
	Key     string
	modKeys modifierKey
}

type KeyPressEvent struct{ KeyboardEvent }
type KeyDownEvent struct{ KeyboardEvent }
type KeyUpEvent struct{ KeyboardEvent }

type modifierKey byte

const (
	modKeyAlt modifierKey = 1 << iota
	modKeyShift
	modKeyCtrl
	modKeyMeta
)

type SendEventMask int

const (
	SendMouseMove SendEventMask = 1 << iota
	SendMouseDown
	SendMouseUp
	SendKeyPress
	SendKeyDown
	SendKeyUp
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
