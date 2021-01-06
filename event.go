// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"errors"
	"math"
)

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

const (
	evMouseMove byte = 1 + iota
	evMouseDown
	evMouseUp
	evKeyPress
	evKeyDown
	evKeyUp
	evClick
	evDblClick
	evAuxClick
	evWheel
	evTouchStart
	evTouchMove
	evTouchEnd
	evTouchCancel
)

func decodeEvent(p []byte) (Event, error) {
	eventType := p[0]
	switch eventType {
	case evMouseMove:
		return MouseMoveEvent{decodeMouseEvent(p)}, nil
	case evMouseDown:
		return MouseDownEvent{decodeMouseEvent(p)}, nil
	case evMouseUp:
		return MouseUpEvent{decodeMouseEvent(p)}, nil
	case evKeyPress:
		return KeyPressEvent{decodeKeyboardEvent(p)}, nil
	case evKeyDown:
		return KeyDownEvent{decodeKeyboardEvent(p)}, nil
	case evKeyUp:
		return KeyUpEvent{decodeKeyboardEvent(p)}, nil
	case evClick:
		return ClickEvent{decodeMouseEvent(p)}, nil
	case evDblClick:
		return DblClickEvent{decodeMouseEvent(p)}, nil
	case evAuxClick:
		return AuxClickEvent{decodeMouseEvent(p)}, nil
	case evWheel:
		return decodeWheelEvent(p), nil
	case evTouchStart:
		return TouchStartEvent{decodeTouchEvent(p)}, nil
	case evTouchMove:
		return TouchMoveEvent{decodeTouchEvent(p)}, nil
	case evTouchEnd:
		return TouchEndEvent{decodeTouchEvent(p)}, nil
	case evTouchCancel:
		return TouchCancelEvent{decodeTouchEvent(p)}, nil
	}
	return nil, errors.New("unknown event type: '" + string(eventType) + "'")
}

func decodeMouseEvent(p []byte) MouseEvent {
	return MouseEvent{
		Buttons:      MouseButtons(p[1]),
		X:            int(byteOrder.Uint32(p[2:])),
		Y:            int(byteOrder.Uint32(p[6:])),
		modifierKeys: modifierKeys(p[10]),
	}
}

func decodeKeyboardEvent(p []byte) KeyboardEvent {
	keyStringLength := int(byteOrder.Uint32(p[2:]))
	return KeyboardEvent{
		Key:          string(p[6 : 6+keyStringLength]),
		modifierKeys: modifierKeys(p[1]),
	}
}

func decodeWheelEvent(p []byte) WheelEvent {
	return WheelEvent{
		MouseEvent: decodeMouseEvent(p),
		DeltaX:     math.Float64frombits(byteOrder.Uint64(p[11:])),
		DeltaY:     math.Float64frombits(byteOrder.Uint64(p[19:])),
		DeltaZ:     math.Float64frombits(byteOrder.Uint64(p[27:])),
		DeltaMode:  DeltaMode(p[35]),
	}
}

func decodeTouchEvent(p []byte) TouchEvent {
	touches, p := decodeTouchList(p[1:])
	changedTouches, p := decodeTouchList(p)
	targetTouches, p := decodeTouchList(p)
	return TouchEvent{
		Touches:        touches,
		ChangedTouches: changedTouches,
		TargetTouches:  targetTouches,
		modifierKeys:   modifierKeys(p[0]),
	}
}

func decodeTouchList(p []byte) (TouchList, []byte) {
	length := p[0]
	p = p[1:]
	list := make(TouchList, length)
	for i := range list {
		var t Touch
		t, p = decodeTouch(p)
		list[i] = t
	}
	return list, p
}

func decodeTouch(p []byte) (Touch, []byte) {
	return Touch{
		Identifier: byteOrder.Uint32(p),
		X:          int(byteOrder.Uint32(p[4:])),
		Y:          int(byteOrder.Uint32(p[8:])),
	}, p[12:]
}
