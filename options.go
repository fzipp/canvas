// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package canvas

import (
	"image/color"
	"time"
)

// Options configure various aspects of the canvas.
// The zero value of this struct is useful out of the box,
// but most users probably want to set at least Width and Height.
type Options struct {
	// Title sets the title of the browser tab/window.
	Title string
	// Width sets the width of the canvas.
	// If Width is not set (i.e. 0) a default value of 300 will be used.
	Width int
	// Height sets the height of the canvas.
	// If Height is not set (i.e. 0) a default value of 150 will be used.
	Height int
	// PageBackground configures the background color
	// of the served HTML page.
	// If PageBackground is not set (i.e. nil) a default value of
	// color.White will be used.
	PageBackground color.Color
	// EnabledEvents enables transmission of the given event types
	// from the client to the server.
	// It is sufficient to list the zero values of the events here.
	// CloseEvent is always implicitly enabled,
	// it doesn't have to be part of this slice.
	EnabledEvents []Event
	// MouseCursorHidden hides the mouse cursor on the canvas.
	MouseCursorHidden bool
	// ContextMenuDisabled disables the context menu on the canvas.
	ContextMenuDisabled bool
	// ScaleToPageWidth scales the canvas to the full horizontal extent
	// of the page in the browser window.
	// This scaling does not change the width within the coordinate system
	// as set by Width.
	ScaleToPageWidth bool
	// ScaleToPageHeight scales the canvas to the full vertical extent
	// of the page in the browser window.
	// This scaling does not change the height within the coordinate system
	// as set by Height.
	ScaleToPageHeight bool
	// ReconnectInterval configures the client to reconnect after
	// the given duration if the WebSocket connection was lost.
	// The client tries to reconnect repeatedly until it is successful.
	// If ReconnectInterval is not set (i.e. 0) the canvas will not try
	// to reconnect if the connection was lost.
	ReconnectInterval time.Duration
}

func (o *Options) applyDefaults() {
	if o.Width == 0 {
		o.Width = 300
	}
	if o.Height == 0 {
		o.Height = 150
	}
	if o.PageBackground == nil {
		o.PageBackground = color.White
	}
}

func (o *Options) eventMask() (mask eventMask) {
	for _, e := range o.EnabledEvents {
		mask |= e.mask()
	}
	return mask
}
