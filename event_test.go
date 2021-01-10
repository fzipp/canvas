// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"reflect"
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
				0b00000011,             // Buttons
				0x00, 0x00, 0x00, 0xc8, // X
				0x00, 0x00, 0x00, 0x96, // Y
				0b00000101, // Modifier keys
			},
			MouseMoveEvent{
				MouseEvent{
					Buttons:      ButtonPrimary | ButtonSecondary,
					X:            200,
					Y:            150,
					modifierKeys: modKeyCtrl | modKeyAlt,
				},
			},
		},
		{
			"KeyDownEvent",
			[]byte{
				0x05,                   // Event type
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
			"WheelEvent",
			[]byte{
				0x0A,                   // Event type
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
			"TouchMoveEvent",
			[]byte{
				0x0C, // Event type

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
					Touches: []Touch{
						{Identifier: 0, X: 3840, Y: 165},
						{Identifier: 1, X: 784, Y: 517},
					},
					ChangedTouches: []Touch{{Identifier: 1, X: 240, Y: 162}},
					TargetTouches:  []Touch{{Identifier: 2, X: 256, Y: 512}},
					modifierKeys:   modKeyAlt | modKeyCtrl,
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
