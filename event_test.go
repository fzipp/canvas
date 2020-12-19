// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import "testing"

func TestMouseEventModifierKeys(t *testing.T) {
	type modifierKeyStates struct {
		altKey   bool
		shiftKey bool
		ctrlKey  bool
		metaKey  bool
	}
	tests := []struct {
		name    string
		modKeys modifierKey
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
			event := MouseEvent{modKeys: tt.modKeys}
			got := modifierKeyStates{
				altKey:   event.AltKey(),
				shiftKey: event.ShiftKey(),
				ctrlKey:  event.CtrlKey(),
				metaKey:  event.MetaKey(),
			}
			if got != tt.want {
				t.Errorf("\ngot : %#v\nwant: %#v", got, tt.want)
			}
		})
	}
}
