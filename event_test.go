// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"reflect"
	"strings"
	"testing"
)

func TestModifierKeys(t *testing.T) {
	type modifierKeyStates struct {
		altKey   bool
		shiftKey bool
		ctrlKey  bool
		metaKey  bool
	}
	tests := []struct {
		name    string
		modKeys modifierKeys
		want    modifierKeyStates
	}{
		{
			name:    "AltKey",
			modKeys: 0b0001,
			want:    modifierKeyStates{altKey: true},
		},
		{
			name:    "ShiftKey",
			modKeys: 0b0010,
			want:    modifierKeyStates{shiftKey: true},
		},
		{
			name:    "CtrlKey",
			modKeys: 0b0100,
			want:    modifierKeyStates{ctrlKey: true},
		},
		{
			name:    "MetaKey",
			modKeys: 0b1000,
			want:    modifierKeyStates{metaKey: true},
		},
		{
			name:    "Multiple modifier keys",
			modKeys: 0b1111,
			want: modifierKeyStates{
				altKey:   true,
				shiftKey: true,
				ctrlKey:  true,
				metaKey:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := modifierKeyStates{
				altKey:   tt.modKeys.AltKey(),
				shiftKey: tt.modKeys.ShiftKey(),
				ctrlKey:  tt.modKeys.CtrlKey(),
				metaKey:  tt.modKeys.MetaKey(),
			}
			if got != tt.want {
				t.Errorf("\ngot : %#v\nwant: %#v", got, tt.want)
			}
		})
	}
}

func TestDecodeEvent(t *testing.T) {
	tests := []struct {
		name string
		p    []byte
		want Event
	}{
		{
			"MouseMoveEvent",
			[]byte{
				0x01,                   // Event type
				0b00000000,             // Buttons
				0x00, 0x00, 0x00, 0xc8, // X
				0x00, 0x00, 0x00, 0x96, // Y
				0b00000101, // Modifier keys
			},
			MouseMoveEvent{
				MouseEvent{
					Buttons:      ButtonNone,
					X:            200,
					Y:            150,
					modifierKeys: modKeyCtrl | modKeyAlt,
				},
			},
		},
		{
			"MouseDownEvent",
			[]byte{
				0x02,                   // Event type
				0b00010000,             // Buttons
				0x00, 0x00, 0x02, 0xfc, // X
				0x00, 0x00, 0x03, 0xff, // Y
				0b00001000, // Modifier keys
			},
			MouseDownEvent{
				MouseEvent{
					Buttons:      Button5th,
					X:            764,
					Y:            1023,
					modifierKeys: modKeyMeta,
				},
			},
		},
		{
			"MouseUpEvent",
			[]byte{
				0x03,                   // Event type
				0b00000001,             // Buttons
				0x00, 0x00, 0x00, 0xf0, // X
				0x00, 0x00, 0x0c, 0x8a, // Y
				0b00000100, // Modifier keys
			},
			MouseUpEvent{
				MouseEvent{
					Buttons:      ButtonPrimary,
					X:            240,
					Y:            3210,
					modifierKeys: modKeyCtrl,
				},
			},
		},
		{
			"KeyDownEvent",
			[]byte{
				0x04,                   // Event type
				0b00001010,             // Modifier keys
				0x00, 0x00, 0x00, 0x09, // len(Key)
				0x41, 0x72, 0x72, 0x6f, 0x77, 0x4c, 0x65, 0x66, 0x74, // Key
			},
			KeyDownEvent{
				KeyboardEvent{
					Key:          "ArrowLeft",
					modifierKeys: modKeyShift | modKeyMeta,
				},
			},
		},
		{
			"KeyUpEvent",
			[]byte{
				0x05,                   // Event type
				0b00000011,             // Modifier keys
				0x00, 0x00, 0x00, 0x0a, // len(Key)
				0x41, 0x72, 0x72, 0x6f, 0x77, 0x52, 0x69, 0x67, 0x68, 0x74, // Key
			},
			KeyUpEvent{
				KeyboardEvent{
					Key:          "ArrowRight",
					modifierKeys: modKeyAlt | modKeyShift,
				},
			},
		},
		{
			"ClickEvent",
			[]byte{
				0x06,                   // Event type
				0b00000010,             // Buttons
				0x00, 0x00, 0x07, 0x81, // X
				0x00, 0x00, 0x02, 0x02, // Y
				0b00000010, // Modifier keys
			},
			ClickEvent{
				MouseEvent{
					Buttons:      ButtonSecondary,
					X:            1921,
					Y:            514,
					modifierKeys: modKeyShift,
				},
			},
		},
		{
			"DblClickEvent",
			[]byte{
				0x07,                   // Event type
				0b00000001,             // Buttons
				0x00, 0x00, 0x0f, 0xc0, // X
				0x00, 0x00, 0x14, 0xdf, // Y
				0b00000011, // Modifier keys
			},
			DblClickEvent{
				MouseEvent{
					Buttons:      ButtonPrimary,
					X:            4032,
					Y:            5343,
					modifierKeys: modKeyAlt | modKeyShift,
				},
			},
		},
		{
			"AuxClickEvent",
			[]byte{
				0x08,                   // Event type
				0b00000001,             // Buttons
				0x00, 0x00, 0x01, 0x41, // X
				0x00, 0x00, 0x02, 0x1f, // Y
				0b00001100, // Modifier keys
			},
			AuxClickEvent{
				MouseEvent{
					Buttons:      ButtonPrimary,
					X:            321,
					Y:            543,
					modifierKeys: modKeyCtrl | modKeyMeta,
				},
			},
		},
		{
			"WheelEvent",
			[]byte{
				0x09,                   // Event type
				0b00001100,             // Buttons
				0x00, 0x00, 0x00, 0x82, // X
				0x00, 0x00, 0x01, 0x9A, // Y
				0b00000010, // Modifier keys

				0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // DeltaX
				0x40, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // DeltaY
				0x40, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // DeltaZ
				0x01, // Delta mode
			},
			WheelEvent{
				MouseEvent: MouseEvent{
					Buttons:      ButtonAuxiliary | Button4th,
					X:            130,
					Y:            410,
					modifierKeys: modKeyShift,
				},
				DeltaX:    10,
				DeltaY:    24,
				DeltaZ:    4,
				DeltaMode: DeltaLine,
			},
		},
		{
			"TouchStartEvent",
			[]byte{
				0x0a,                   // Event type
				0x01,                   // len(Touches)
				0x00, 0x00, 0x00, 0x00, // Touches[0].Identifier
				0x00, 0x00, 0x01, 0x54, // Touches[0].X
				0x00, 0x00, 0x00, 0xd2, // Touches[0].Y
				0x00,       // len(ChangedTouches)
				0x00,       // len(TargetTouches)
				0b00001000, // Modifier keys
			},
			TouchStartEvent{
				TouchEvent{
					Touches: TouchList{
						{Identifier: 0, X: 340, Y: 210},
					},
					ChangedTouches: TouchList{},
					TargetTouches:  TouchList{},
					modifierKeys:   modKeyMeta,
				},
			},
		},
		{
			"TouchMoveEvent",
			[]byte{
				0x0b, // Event type

				0x02,                   // len(Touches)
				0x00, 0x00, 0x00, 0x00, // Touches[0].Identifier
				0x00, 0x00, 0x0f, 0x00, // Touches[0].X
				0x00, 0x00, 0x00, 0xa5, // Touches[0].Y
				0x00, 0x00, 0x00, 0x01, // Touches[1].Identifier
				0x00, 0x00, 0x03, 0x10, // Touches[1].X
				0x00, 0x00, 0x02, 0x05, // Touches[1].Y

				0x01,                   // len(ChangedTouches)
				0x00, 0x00, 0x00, 0x01, // ChangedTouches[0].Identifier
				0x00, 0x00, 0x00, 0xf0, // ChangedTouches[0].X
				0x00, 0x00, 0x00, 0xa2, // ChangedTouches[0].Y

				0x01,                   // len(TargetTouches)
				0x00, 0x00, 0x00, 0x02, // TargetTouches[0].Identifier
				0x00, 0x00, 0x01, 0x00, // TargetTouches[0].X
				0x00, 0x00, 0x02, 0x00, // TargetTouches[0].Y

				0b00000101, // Modifier keys
			},
			TouchMoveEvent{
				TouchEvent{
					Touches: TouchList{
						{Identifier: 0, X: 3840, Y: 165},
						{Identifier: 1, X: 784, Y: 517},
					},
					ChangedTouches: TouchList{{Identifier: 1, X: 240, Y: 162}},
					TargetTouches:  TouchList{{Identifier: 2, X: 256, Y: 512}},
					modifierKeys:   modKeyAlt | modKeyCtrl,
				},
			},
		},
		{
			"TouchEndEvent",
			[]byte{
				0x0c,                   // Event type
				0x00,                   // len(Touches)
				0x01,                   // len(ChangedTouches)
				0x00, 0x00, 0x00, 0x00, // Touches[0].Identifier
				0x00, 0x00, 0x01, 0x54, // Touches[0].X
				0x00, 0x00, 0x00, 0xd2, // Touches[0].Y
				0x00,       // len(TargetTouches)
				0b00000001, // Modifier keys
			},
			TouchEndEvent{
				TouchEvent{
					Touches: TouchList{},
					ChangedTouches: TouchList{
						{Identifier: 0, X: 340, Y: 210},
					},
					TargetTouches: TouchList{},
					modifierKeys:  modKeyAlt,
				},
			},
		},
		{
			"TouchCancelEvent",
			[]byte{
				0x0d,                   // Event type
				0x00,                   // len(Touches)
				0x00,                   // len(ChangedTouches)
				0x01,                   // len(TargetTouches)
				0x00, 0x00, 0x00, 0x00, // Touches[0].Identifier
				0x00, 0x00, 0x01, 0x54, // Touches[0].X
				0x00, 0x00, 0x00, 0xd2, // Touches[0].Y
				0b00000011, // Modifier keys
			},
			TouchCancelEvent{
				TouchEvent{
					Touches:        TouchList{},
					ChangedTouches: TouchList{},
					TargetTouches: TouchList{
						{Identifier: 0, X: 340, Y: 210},
					},
					modifierKeys: modKeyAlt | modKeyShift,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeEvent(tt.p)
			if err != nil {
				t.Errorf("did not expect error, but got error: %s", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot : %#v\nwant: %#v", got, tt.want)
			}
		})
	}
}

func TestUnsupportedEventType(t *testing.T) {
	tests := []byte{
		0x00,
		0xfa,
		0xfb,
		0xfc,
	}
	wantErrorMessage := "unknown event type: "
	for _, tt := range tests {
		got, err := decodeEvent([]byte{tt})
		if err == nil {
			t.Errorf("expected error, but got none")
			return
		}
		if wantType, ok := err.(errUnknownEventType); !ok {
			t.Errorf("expected %T error, but got: %#v", wantType, err)
		}
		if !strings.HasPrefix(err.Error(), wantErrorMessage) {
			t.Errorf("expected %q error message, but got: %q", wantErrorMessage, err)
		}
		if got != nil {
			t.Errorf("expected nil event, but got: %#v", got)
		}
	}
}

func TestEventDataTooShort(t *testing.T) {
	tests := []struct {
		p []byte
	}{
		{[]byte{0x01}},
		{[]byte{0x02, 0x00, 0x00, 0x00}},
		{[]byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}
	for _, tt := range tests {
		got, err := decodeEvent(tt.p)
		if err == nil {
			t.Errorf("expected error, but got none")
			return
		}
		if wantType, ok := err.(errDataTooShort); !ok {
			t.Errorf("expected %T error, but got: %#v", wantType, err)
		}
		if got != nil {
			t.Errorf("expected nil event, but got: %#v", got)
		}
	}
}

func TestEventMask(t *testing.T) {
	tests := []struct {
		event Event
		want  eventMask
	}{
		{CloseEvent{}, 0},
		{MouseEvent{}, 0b11100111},
		{MouseMoveEvent{}, 0b00000001},
		{MouseDownEvent{}, 0b00000010},
		{MouseUpEvent{}, 0b00000100},
		{ClickEvent{}, 0b00100000},
		{DblClickEvent{}, 0b01000000},
		{AuxClickEvent{}, 0b10000000},
		{WheelEvent{}, 0b100000000},
		{KeyboardEvent{}, 0b00011000},
		{KeyDownEvent{}, 0b00001000},
		{KeyUpEvent{}, 0b00010000},
		{TouchEvent{}, 0b1111000000000},
		{TouchStartEvent{}, 0b0001000000000},
		{TouchMoveEvent{}, 0b0010000000000},
		{TouchEndEvent{}, 0b0100000000000},
		{TouchCancelEvent{}, 0b1000000000000},
	}
	for _, tt := range tests {
		got := tt.event.mask()
		if got != tt.want {
			t.Errorf("Event mask for %T; got: %#b, want: %#b", tt.event, got, tt.want)
		}
	}
}
