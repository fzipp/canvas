// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The API doc comments are based on the MDN Web Docs for the [Canvas API]
// by Mozilla Contributors and are licensed under [CC-BY-SA 2.5].
//
// [Canvas API]: https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D
// [CC-BY-SA 2.5]: https://creativecommons.org/licenses/by-sa/2.5/

package canvas

import "fmt"

// Event is an interface implemented by all event subtypes. Events can be
// received from the channel returned by Context.Events. Use a type switch
// to distinguish between different event types.
type Event interface {
	mask() eventMask
}

// The CloseEvent is fired when the WebSocket connection to the client is
// closed.
// It is not necessary to enable the CloseEvent with Options.EnabledEvents,
// it is always enabled.
// Animation loops should handle the CloseEvent to quit the loop.
type CloseEvent struct{}

func (e CloseEvent) mask() eventMask { return 0 }

// MouseEvent represents events that occur due to the user interacting with a
// pointing device (such as a mouse).
type MouseEvent struct {
	// Buttons encodes the buttons being depressed (if any) when the mouse
	// event was fired.
	Buttons MouseButtons
	// The X coordinate of the mouse pointer.
	X int
	// The Y coordinate of the mouse pointer.
	Y int
	// Mod describes the modifier keys pressed during the event.
	Mod ModifierKeys
}

func (e MouseEvent) mask() eventMask {
	return maskMouseMove | maskMouseUp | maskMouseDown | maskClick | maskDblClick | maskAuxClick
}

// The MouseMoveEvent is fired when a pointing device (usually a mouse) is
// moved.
type MouseMoveEvent struct{ MouseEvent }

func (e MouseMoveEvent) mask() eventMask { return maskMouseMove }

// The MouseDownEvent is fired when a pointing device button is pressed.
//
// Note: This differs from the ClickEvent in that click is fired after a full
// click action occurs; that is, the mouse button is pressed and released
// while the pointer remains inside the canvas. MouseDownEvent is fired
// the moment the button is initially pressed.
type MouseDownEvent struct{ MouseEvent }

func (e MouseDownEvent) mask() eventMask { return maskMouseDown }

// The MouseUpEvent is fired when a button on a pointing device (such as a
// mouse or trackpad) is released. It is the counterpoint to the
// MouseDownEvent.
type MouseUpEvent struct{ MouseEvent }

func (e MouseUpEvent) mask() eventMask { return maskMouseUp }

// The ClickEvent is fired when a pointing device button (such as a mouse's
// primary mouse button) is both pressed and released while the pointer is
// located inside the canvas.
type ClickEvent struct{ MouseEvent }

func (e ClickEvent) mask() eventMask { return maskClick }

// The DblClickEvent is fired when a pointing device button (such as a mouse's
// primary button) is double-clicked; that is, when it's rapidly clicked twice
// on the canvas within a very short span of time.
//
// DblClickEvent fires after two ClickEvents (and by extension, after two pairs
// of MouseDownEvents and MouseUpEvents).
type DblClickEvent struct{ MouseEvent }

func (e DblClickEvent) mask() eventMask { return maskDblClick }

// The AuxClickEvent is fired when a non-primary pointing device button (any
// mouse button other than the primary—usually leftmost—button) has been
// pressed and released both within the same element.
type AuxClickEvent struct{ MouseEvent }

func (e AuxClickEvent) mask() eventMask { return maskAuxClick }

// The WheelEvent is fired due to the user moving a mouse wheel or similar
// input device.
type WheelEvent struct {
	MouseEvent
	// DeltaX represents the horizontal scroll amount.
	DeltaX float64
	// DeltaY represents the vertical scroll amount.
	DeltaY float64
	// DeltaZ represents the scroll amount for the z-axis.
	DeltaZ float64
	// DeltaMode represents the unit of the delta values' scroll amount.
	DeltaMode DeltaMode
}

func (e WheelEvent) mask() eventMask {
	return maskWheel
}

// DeltaMode represents the unit of the delta values' scroll amount.
type DeltaMode byte

const (
	// DeltaPixel means the delta values are specified in pixels.
	DeltaPixel DeltaMode = iota
	// DeltaLine means the delta values are specified in lines.
	DeltaLine
	// DeltaPage means the delta values are specified in pages.
	DeltaPage
)

// KeyboardEvent objects describe a user interaction with the keyboard; each
// event describes a single interaction between the user and a key (or
// combination of a key with modifier keys) on the keyboard.
type KeyboardEvent struct {
	// Key represents the key value of the key represented by the event.
	Key string
	// Mod describes the modifier keys pressed during the event.
	Mod ModifierKeys
}

func (e KeyboardEvent) mask() eventMask {
	return maskKeyDown | maskKeyUp
}

// The KeyDownEvent is fired when a key is pressed.
type KeyDownEvent struct{ KeyboardEvent }

func (e KeyDownEvent) mask() eventMask { return maskKeyDown }

// The KeyUpEvent is fired when a key is released.
type KeyUpEvent struct{ KeyboardEvent }

func (e KeyUpEvent) mask() eventMask { return maskKeyUp }

// The TouchEvent is fired when the state of contacts with a touch-sensitive
// surface changes. This surface can be a touch screen or trackpad, for
// example. The event can describe one or more points of contact with the
// screen and includes support for detecting movement, addition and removal of
// contact points, and so forth.
//
// Touches are represented by the Touch object; each touch is described by a
// position, size and shape, amount of pressure, and target element. Lists of
// touches are represented by TouchList objects.
type TouchEvent struct {
	// Touches is a TouchList of all the Touch objects representing all current
	// points of contact with the surface, regardless of target or changed
	// status.
	Touches TouchList
	// ChangedTouches is a TouchList of all the Touch objects representing
	// individual points of contact whose states changed between the previous
	// touch event and this one.
	ChangedTouches TouchList
	// TargetTouches is a TouchList of all the Touch objects that are both
	// currently in contact with the touch surface and were also started on the
	// same element that is the target of the event.
	TargetTouches TouchList
	// Mod describes the modifier keys pressed during the event.
	Mod ModifierKeys
}

func (e TouchEvent) mask() eventMask {
	return maskTouchStart | maskTouchMove | maskTouchEnd | maskTouchCancel
}

// TouchList represents a list of contact points on a touch surface. For
// example, if the user has three fingers on the touch surface (such as a
// screen or trackpad), the corresponding TouchList object would have one
// Touch object for each finger, for a total of three entries.
type TouchList []Touch

// Touch represents a single contact point on a touch-sensitive device.
// The contact point is commonly a finger or stylus and the device may be a
// touchscreen or trackpad.
type Touch struct {
	// Identifier is a unique identifier for this Touch object. A given touch
	// point (say, by a finger) will have the same identifier for the duration
	// of its movement around the surface. This lets you ensure that you're
	// tracking the same touch all the time.
	Identifier uint32
	// The X coordinate of the touch point.
	X int
	// The Y coordinate of the touch point.
	Y int
}

// The TouchStartEvent is fired when one or more touch points are placed on
// the touch surface.
type TouchStartEvent struct{ TouchEvent }

func (e TouchStartEvent) mask() eventMask { return maskTouchStart }

// The TouchMoveEvent is fired when one or more touch points are moved along
// the touch surface.
type TouchMoveEvent struct{ TouchEvent }

func (e TouchMoveEvent) mask() eventMask { return maskTouchMove }

// The TouchEndEvent is fired when one or more touch points are removed from
// the touch surface.
type TouchEndEvent struct{ TouchEvent }

func (e TouchEndEvent) mask() eventMask { return maskTouchEnd }

// The TouchCancelEvent is fired when one or more touch points have been
// disrupted in an implementation-specific manner (for example, too many touch
// points are created).
type TouchCancelEvent struct{ TouchEvent }

func (e TouchCancelEvent) mask() eventMask { return maskTouchCancel }

// ModifierKeys describes the modifier keys (Alt, Shift, Ctrl, Meta) pressed
// during an event.
type ModifierKeys byte

const (
	modKeyAlt ModifierKeys = 1 << iota
	modKeyShift
	modKeyCtrl
	modKeyMeta
)

// AltKey returns true if the Alt (Option or ⌥ on OS X) key was active when
// the event was generated.
func (m ModifierKeys) AltKey() bool {
	return m.isPressed(modKeyAlt)
}

// ShiftKey returns true if the Shift key was active when the event was
// generated.
func (m ModifierKeys) ShiftKey() bool {
	return m.isPressed(modKeyShift)
}

// CtrlKey returns true if the Ctrl key was active when the event was
// generated.
func (m ModifierKeys) CtrlKey() bool {
	return m.isPressed(modKeyCtrl)
}

// MetaKey returns true if the Meta key (on Mac keyboards, the ⌘ Command key;
// on Windows keyboards, the Windows key (⊞)) was active when the event
// was generated.
func (m ModifierKeys) MetaKey() bool {
	return m.isPressed(modKeyMeta)
}

func (m ModifierKeys) isPressed(k ModifierKeys) bool {
	return m&k != 0
}

type eventMask int

const (
	maskMouseMove eventMask = 1 << iota
	maskMouseDown
	maskMouseUp
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

// MouseButtons is a number representing one or more buttons. For more than
// one button pressed simultaneously, the values are combined (e.g., 3 is
// ButtonPrimary + ButtonSecondary).
type MouseButtons int

const (
	// ButtonPrimary is the primary button (usually the left button).
	ButtonPrimary MouseButtons = 1 << iota
	// ButtonSecondary is the secondary button (usually the right button).
	ButtonSecondary
	// ButtonAuxiliary is the auxiliary button (usually the mouse wheel button
	// or middle button)
	ButtonAuxiliary
	// Button4th is the 4th button (typically the "Browser Back" button).
	Button4th
	// Button5th is the 5th button (typically the "Browser Forward" button).
	Button5th
	// ButtonNone stands for no button or un-initialized.
	ButtonNone MouseButtons = 0
)

const (
	evMouseMove byte = 1 + iota
	evMouseDown
	evMouseUp
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
	return nil, errUnknownEventType{unknownType: eventType}
}

func decodeMouseEvent(buf *buffer) MouseEvent {
	return MouseEvent{
		Buttons: MouseButtons(buf.readByte()),
		X:       int(buf.readUint32()),
		Y:       int(buf.readUint32()),
		Mod:     ModifierKeys(buf.readByte()),
	}
}

func decodeKeyboardEvent(buf *buffer) KeyboardEvent {
	return KeyboardEvent{
		Mod: ModifierKeys(buf.readByte()),
		Key: buf.readString(),
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
		Mod:            ModifierKeys(buf.readByte()),
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

type errUnknownEventType struct {
	unknownType byte
}

func (err errUnknownEventType) Error() string {
	return fmt.Sprintf("unknown event type: %#x", err.unknownType)
}
