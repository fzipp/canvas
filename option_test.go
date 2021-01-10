// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"testing"
	"time"
)

func TestConfigFrom(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
		want    config
	}{
		{
			"Title",
			[]Option{Title("Test title")},
			config{title: "Test title"},
		},
		{
			"Size",
			[]Option{Size(800, 600)},
			config{width: 800, height: 600},
		},
		{
			"EnableEvents",
			[]Option{EnableEvents(MouseMoveEvent{}, MouseDownEvent{}, KeyDownEvent{})},
			config{eventMask: 0b10011},
		},
		{
			"DisableCursor",
			[]Option{DisableCursor()},
			config{cursorDisabled: true},
		},
		{
			"DisableContextMenu",
			[]Option{DisableContextMenu()},
			config{contextMenuDisabled: true},
		},
		{
			"FullPage",
			[]Option{FullPage()},
			config{fullPage: true},
		},
		{
			"Reconnect",
			[]Option{Reconnect(2 * time.Second)},
			config{reconnectInterval: 2 * time.Second},
		},
		{
			"Multiple options",
			[]Option{
				Title("hello, world"),
				Size(320, 200),
				EnableEvents(MouseMoveEvent{}, MouseDownEvent{}, MouseUpEvent{}),
				DisableCursor(),
				DisableContextMenu(),
				FullPage(),
				Reconnect(1500 * time.Millisecond),
			},
			config{
				title:               "hello, world",
				width:               320,
				height:              200,
				eventMask:           0b0111,
				cursorDisabled:      true,
				contextMenuDisabled: true,
				fullPage:            true,
				reconnectInterval:   1500 * time.Millisecond,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := configFrom(tt.options)
			if got != tt.want {
				t.Errorf("\ngot : %#v\nwant: %#v", got, tt.want)
			}
		})
	}
}
