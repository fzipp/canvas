// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image/color"
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
			config{
				title:           "Test title",
				width:           300,
				height:          150,
				backgroundColor: color.White,
			},
		},
		{
			"Size",
			[]Option{Size(800, 600)},
			config{
				width:           800,
				height:          600,
				backgroundColor: color.White,
			},
		},
		{
			"EnableEvents",
			[]Option{EnableEvents(MouseMoveEvent{}, MouseDownEvent{}, KeyDownEvent{})},
			config{
				eventMask:       0b1011,
				width:           300,
				height:          150,
				backgroundColor: color.White,
			},
		},
		{
			"DisableCursor",
			[]Option{DisableCursor()},
			config{
				cursorDisabled:  true,
				width:           300,
				height:          150,
				backgroundColor: color.White,
			},
		},
		{
			"DisableContextMenu",
			[]Option{DisableContextMenu()},
			config{
				contextMenuDisabled: true,
				width:               300,
				height:              150,
				backgroundColor:     color.White,
			},
		},
		{
			"FullPage",
			[]Option{FullPage()},
			config{
				fullPage:        true,
				width:           300,
				height:          150,
				backgroundColor: color.White,
			},
		},
		{
			"Reconnect",
			[]Option{Reconnect(2 * time.Second)},
			config{
				reconnectInterval: 2 * time.Second,
				width:             300,
				height:            150,
				backgroundColor:   color.White,
			},
		},
		{
			"Multiple options",
			[]Option{
				Title("hello, world"),
				Size(320, 200),
				BackgroundColor(color.Black),
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
				backgroundColor:     color.Black,
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
