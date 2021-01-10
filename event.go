// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import "fmt"

type Event interface {
	mask() eventMask
}

type CloseEvent struct {}

func (e CloseEvent) mask() eventMask { return 0 }

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
	buf := &buffer{bytes: p}
	event, err := decodeEventBuf(buf)
	if buf.error != nil {
		return nil, buf.error
	}
	return event, err
}

func decodeEventBuf(buf *buffer) (Event, error) {
	eventType := buf.readByte()
	switch eventType {
	case evMouseMove:
		return MouseMoveEvent{decodeMouseEvent(buf)}, nil
	case evMouseDown:
		return MouseDownEvent{decodeMouseEvent(buf)}, nil
	case evMouseUp:
		return MouseUpEvent{decodeMouseEvent(buf)}, nil
	case evKeyPress:
		return KeyPressEvent{decodeKeyboardEvent(buf)}, nil
	case evKeyDown:
		return KeyDownEvent{decodeKeyboardEvent(buf)}, nil
	case evKeyUp:
		return KeyUpEvent{decodeKeyboardEvent(buf)}, nil
	case evClick:
		return ClickEvent{decodeMouseEvent(buf)}, nil
	case evDblClick:
		return DblClickEvent{decodeMouseEvent(buf)}, nil
	case evAuxClick:
		return AuxClickEvent{decodeMouseEvent(buf)}, nil
	case evWheel:
		return decodeWheelEvent(buf), nil
	case evTouchStart:
		return TouchStartEvent{decodeTouchEvent(buf)}, nil
	case evTouchMove:
		return TouchMoveEvent{decodeTouchEvent(buf)}, nil
	case evTouchEnd:
		return TouchEndEvent{decodeTouchEvent(buf)}, nil
	case evTouchCancel:
		return TouchCancelEvent{decodeTouchEvent(buf)}, nil
	}
	return nil, fmt.Errorf("unknown event type: %#x", eventType)
}

func decodeMouseEvent(buf *buffer) MouseEvent {
	return MouseEvent{
		Buttons:      MouseButtons(buf.readByte()),
		X:            int(buf.readUint32()),
		Y:            int(buf.readUint32()),
		modifierKeys: modifierKeys(buf.readByte()),
	}
}

func decodeKeyboardEvent(buf *buffer) KeyboardEvent {
	return KeyboardEvent{
		modifierKeys: modifierKeys(buf.readByte()),
		Key:          buf.readString(),
	}
}

func decodeWheelEvent(buf *buffer) WheelEvent {
	return WheelEvent{
		MouseEvent: decodeMouseEvent(buf),
		DeltaX:     buf.readFloat64(),
		DeltaY:     buf.readFloat64(),
		DeltaZ:     buf.readFloat64(),
		DeltaMode:  DeltaMode(buf.readByte()),
	}
}

func decodeTouchEvent(buf *buffer) TouchEvent {
	return TouchEvent{
		Touches:        decodeTouchList(buf),
		ChangedTouches: decodeTouchList(buf),
		TargetTouches:  decodeTouchList(buf),
		modifierKeys:   modifierKeys(buf.readByte()),
	}
}

func decodeTouchList(buf *buffer) TouchList {
	length := buf.readByte()
	list := make(TouchList, length)
	for i := range list {
		list[i] = decodeTouch(buf)
	}
	return list
}

func decodeTouch(buf *buffer) Touch {
	return Touch{
		Identifier: buf.readUint32(),
		X:          int(buf.readUint32()),
		Y:          int(buf.readUint32()),
	}
}
