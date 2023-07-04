// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image/color"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOptionsApplyDefaults(t *testing.T) {
	tests := []struct {
		name string
		opts *Options
		want *Options
	}{
		{
			"empty options",
			&Options{},
			&Options{
				Width:          300,
				Height:         150,
				PageBackground: color.White,
			},
		},
		{
			"width and height given",
			&Options{
				Width:  800,
				Height: 600,
			},
			&Options{
				Width:          800,
				Height:         600,
				PageBackground: color.White,
			},
		},
		{
			"background color given",
			&Options{
				PageBackground: color.Black,
			},
			&Options{
				Width:          300,
				Height:         150,
				PageBackground: color.Black,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := *tt.opts
			got.applyDefaults()
			if diff := cmp.Diff(tt.want, &got, cmp.AllowUnexported(Options{})); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOptionsEventMask(t *testing.T) {
	tests := []struct {
		name string
		opts *Options
		want eventMask
	}{
		{
			"empty options",
			&Options{},
			0,
		},
		{
			"multiple events",
			&Options{
				EnabledEvents: []Event{
					KeyUpEvent{},
					MouseMoveEvent{},
					TouchStartEvent{},
				},
			},
			0b1000010001,
		},
		{
			"keyboard events",
			&Options{
				EnabledEvents: []Event{
					KeyboardEvent{},
				},
			},
			0b00011000,
		},
		{
			"mouse events",
			&Options{
				EnabledEvents: []Event{
					MouseEvent{},
				},
			},
			0b11100111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.eventMask()
			if got != tt.want {
				t.Errorf("opts.EnabledEvents = %#v\nopts.eventMask() = %#b, want: %#b",
					tt.opts.EnabledEvents, got, tt.want)
			}
		})
	}
}
